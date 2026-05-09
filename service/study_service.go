package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"
)

type studyResponse struct {
	Code int       `json:"code"`
	Data studyData `json:"data"`
}

type studyData struct {
	App           studyApp      `json:"app"`
	Profile       studyProfile  `json:"profile"`
	Today         studyToday    `json:"today"`
	RecentRecords []studyRecord `json:"recentRecords"`
	Notes         []studyNote   `json:"notes"`
	WeeklyFocus   []string      `json:"weeklyFocus"`
	Footer        studyFooter   `json:"footer"`
}

type studyApp struct {
	Name    string `json:"name"`
	Tagline string `json:"tagline"`
}

type studyProfile struct {
	UserName string `json:"userName"`
	Streak   int    `json:"streak"`
	Total    int    `json:"total"`
	Goal     string `json:"goal"`
}

type studyToday struct {
	Date       string      `json:"date"`
	Motto      string      `json:"motto"`
	Completion string      `json:"completion"`
	Tasks      []studyTask `json:"tasks"`
}

type studyTask struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	Duration string `json:"duration"`
	Status   string `json:"status"`
}

type studyRecord struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Category  string `json:"category"`
	Summary   string `json:"summary"`
	CreatedAt string `json:"createdAt"`
}

type studyNote struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Tag       string `json:"tag"`
	CreatedAt string `json:"createdAt"`
}

type studyFooter struct {
	Nav []string `json:"nav"`
}

type createRecordRequest struct {
	Title    string `json:"title"`
	Category string `json:"category"`
	Summary  string `json:"summary"`
	Note     string `json:"note"`
}

type createTaskRequest struct {
	Title    string `json:"title"`
	Duration string `json:"duration"`
	Status   string `json:"status"`
}

type updateTaskStatusRequest struct {
	Status string `json:"status"`
}

type updateRecordRequest struct {
	Title    string `json:"title"`
	Category string `json:"category"`
	Summary  string `json:"summary"`
}

type updateNoteRequest struct {
	Title   string `json:"title"`
	Tag     string `json:"tag"`
	Content string `json:"content"`
}

type studyRepository interface {
	ListTasks() ([]studyTask, error)
	CreateTask(createTaskRequest) (studyTask, error)
	UpdateTaskStatus(uint, string) error
	DeleteTask(uint) error
	ListRecords(int) ([]studyRecord, error)
	CreateRecord(createRecordRequest) (studyRecord, studyNote, error)
	UpdateRecord(uint, updateRecordRequest) (studyRecord, error)
	DeleteRecord(uint) error
	ListNotes(int) ([]studyNote, error)
	UpdateNote(uint, updateNoteRequest) (studyNote, error)
	DeleteNote(uint) error
}

var studyRepo studyRepository = newMemoryStudyRepository()

func SetStudyRepository(repo studyRepository) {
	if repo != nil {
		studyRepo = repo
	}
}

type dbStudyRepository struct {
	store dao.StudyStore
}

func NewDBStudyRepository(store dao.StudyStore) studyRepository {
	repo := &dbStudyRepository{store: store}
	repo.ensureSeedData()
	return repo
}

func (r *dbStudyRepository) ListTasks() ([]studyTask, error) {
	items, err := r.store.ListTasks()
	if err != nil {
		return nil, err
	}

	return mapTasks(items), nil
}

func (r *dbStudyRepository) CreateTask(req createTaskRequest) (studyTask, error) {
	task := &model.StudyTaskModel{
		Title:    strings.TrimSpace(req.Title),
		Duration: strings.TrimSpace(req.Duration),
		Status:   normalizeTaskStatus(req.Status),
		Sort:     int(time.Now().Unix()),
	}
	if task.Title == "" || task.Duration == "" {
		return studyTask{}, fmt.Errorf("任务标题和时长不能为空")
	}

	if err := r.store.CreateTask(task); err != nil {
		return studyTask{}, err
	}

	return mapTask(*task), nil
}

func (r *dbStudyRepository) UpdateTaskStatus(id uint, status string) error {
	return r.store.UpdateTaskStatus(id, normalizeTaskStatus(status))
}

func (r *dbStudyRepository) DeleteTask(id uint) error {
	return r.store.DeleteTask(id)
}

func (r *dbStudyRepository) ListRecords(limit int) ([]studyRecord, error) {
	items, err := r.store.ListRecords(limit)
	if err != nil {
		return nil, err
	}

	return mapRecords(items), nil
}

func (r *dbStudyRepository) CreateRecord(req createRecordRequest) (studyRecord, studyNote, error) {
	title := strings.TrimSpace(req.Title)
	category := strings.TrimSpace(req.Category)
	summary := strings.TrimSpace(req.Summary)
	noteContent := strings.TrimSpace(req.Note)
	if title == "" || category == "" || summary == "" || noteContent == "" {
		return studyRecord{}, studyNote{}, fmt.Errorf("标题、分类、记录和笔记不能为空")
	}

	recordModel := &model.StudyRecordModel{
		Title:    title,
		Category: category,
		Summary:  summary,
	}
	if err := r.store.CreateRecord(recordModel); err != nil {
		return studyRecord{}, studyNote{}, err
	}

	noteModel := &model.StudyNoteModel{
		Title:   title + " 笔记",
		Content: noteContent,
		Tag:     category,
	}
	if err := r.store.CreateNote(noteModel); err != nil {
		return studyRecord{}, studyNote{}, err
	}

	return mapRecord(*recordModel), mapNote(*noteModel), nil
}

func (r *dbStudyRepository) ListNotes(limit int) ([]studyNote, error) {
	items, err := r.store.ListNotes(limit)
	if err != nil {
		return nil, err
	}

	return mapNotes(items), nil
}

func (r *dbStudyRepository) UpdateRecord(id uint, req updateRecordRequest) (studyRecord, error) {
	title := strings.TrimSpace(req.Title)
	category := strings.TrimSpace(req.Category)
	summary := strings.TrimSpace(req.Summary)
	if title == "" || category == "" || summary == "" {
		return studyRecord{}, fmt.Errorf("标题、分类和记录不能为空")
	}

	recordModel := &model.StudyRecordModel{
		ID:       id,
		Title:    title,
		Category: category,
		Summary:  summary,
	}
	if err := r.store.UpdateRecord(recordModel); err != nil {
		return studyRecord{}, err
	}

	items, err := r.store.ListRecords(0)
	if err != nil {
		return studyRecord{}, err
	}
	for _, item := range items {
		if item.ID == id {
			return mapRecord(item), nil
		}
	}
	return studyRecord{}, fmt.Errorf("记录不存在")
}

func (r *dbStudyRepository) DeleteRecord(id uint) error {
	return r.store.DeleteRecord(id)
}

func (r *dbStudyRepository) UpdateNote(id uint, req updateNoteRequest) (studyNote, error) {
	title := strings.TrimSpace(req.Title)
	tag := strings.TrimSpace(req.Tag)
	content := strings.TrimSpace(req.Content)
	if title == "" || tag == "" || content == "" {
		return studyNote{}, fmt.Errorf("标题、标签和笔记不能为空")
	}

	noteModel := &model.StudyNoteModel{
		ID:      id,
		Title:   title,
		Tag:     tag,
		Content: content,
	}
	if err := r.store.UpdateNote(noteModel); err != nil {
		return studyNote{}, err
	}

	items, err := r.store.ListNotes(0)
	if err != nil {
		return studyNote{}, err
	}
	for _, item := range items {
		if item.ID == id {
			return mapNote(item), nil
		}
	}
	return studyNote{}, fmt.Errorf("笔记不存在")
}

func (r *dbStudyRepository) DeleteNote(id uint) error {
	return r.store.DeleteNote(id)
}

func (r *dbStudyRepository) ensureSeedData() {
	tasks, err := r.store.ListTasks()
	if err == nil && len(tasks) == 0 {
		_ = r.store.CreateTask(&model.StudyTaskModel{Title: "背 20 个英语单词", Duration: "25 分钟", Status: "已完成", Sort: 1})
		_ = r.store.CreateTask(&model.StudyTaskModel{Title: "复盘 Go 路由与接口结构", Duration: "40 分钟", Status: "进行中", Sort: 2})
		_ = r.store.CreateTask(&model.StudyTaskModel{Title: "整理错题与复习笔记", Duration: "20 分钟", Status: "待开始", Sort: 3})
	}

	records, err := r.store.ListRecords(1)
	if err == nil && len(records) == 0 {
		_ = r.store.CreateRecord(&model.StudyRecordModel{Title: "Go HTTP 服务启动流程", Category: "后端", Summary: "把路由注册、端口配置和静态页入口重新梳理了一遍。"})
		_ = r.store.CreateRecord(&model.StudyRecordModel{Title: "英语单词 Day 12", Category: "英语", Summary: "今天重点复习 travel、schedule、focus 这组高频词。"})
	}

	notes, err := r.store.ListNotes(1)
	if err == nil && len(notes) == 0 {
		_ = r.store.CreateNote(&model.StudyNoteModel{Title: "今天的复盘", Content: "学习任务不要排太满，控制在 3 项内更容易完成。", Tag: "方法"})
		_ = r.store.CreateNote(&model.StudyNoteModel{Title: "接口设计提醒", Content: "前端只需要一个 dashboard 接口时，先把返回结构稳定下来。", Tag: "开发"})
	}
}

type memoryStudyRepository struct {
	mu      sync.RWMutex
	nextID  uint
	tasks   []studyTask
	records []studyRecord
	notes   []studyNote
}

func newMemoryStudyRepository() *memoryStudyRepository {
	now := time.Now()
	return &memoryStudyRepository{
		nextID: 9,
		tasks: []studyTask{
			{ID: 1, Title: "背 20 个英语单词", Duration: "25 分钟", Status: "已完成"},
			{ID: 2, Title: "复盘 Go 路由与接口结构", Duration: "40 分钟", Status: "进行中"},
			{ID: 3, Title: "整理错题与复习笔记", Duration: "20 分钟", Status: "待开始"},
		},
		records: []studyRecord{
			{ID: 4, Title: "Go HTTP 服务启动流程", Category: "后端", Summary: "把路由注册、端口配置和静态页入口重新梳理了一遍。", CreatedAt: now.Add(-2 * time.Hour).Format(timeLayout)},
			{ID: 5, Title: "英语单词 Day 12", Category: "英语", Summary: "今天重点复习 travel、schedule、focus 这组高频词。", CreatedAt: now.Add(-26 * time.Hour).Format(timeLayout)},
			{ID: 6, Title: "算法错题整理", Category: "刷题", Summary: "双指针边界判断还是容易漏，明天继续做两题巩固。", CreatedAt: now.Add(-50 * time.Hour).Format(timeLayout)},
		},
		notes: []studyNote{
			{ID: 7, Title: "今天的复盘", Content: "学习任务不要排太满，控制在 3 项内更容易完成。", Tag: "方法", CreatedAt: now.Add(-90 * time.Minute).Format(timeLayout)},
			{ID: 8, Title: "接口设计提醒", Content: "前端只需要一个 dashboard 接口时，先把返回结构稳定下来。", Tag: "开发", CreatedAt: now.Add(-5 * time.Hour).Format(timeLayout)},
		},
	}
}

func (r *memoryStudyRepository) ListTasks() ([]studyTask, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return append([]studyTask(nil), r.tasks...), nil
}

func (r *memoryStudyRepository) CreateTask(req createTaskRequest) (studyTask, error) {
	title := strings.TrimSpace(req.Title)
	duration := strings.TrimSpace(req.Duration)
	if title == "" || duration == "" {
		return studyTask{}, fmt.Errorf("任务标题和时长不能为空")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	task := studyTask{
		ID:       r.nextID,
		Title:    title,
		Duration: duration,
		Status:   normalizeTaskStatus(req.Status),
	}
	r.nextID++
	r.tasks = append(r.tasks, task)
	return task, nil
}

func (r *memoryStudyRepository) UpdateTaskStatus(id uint, status string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := range r.tasks {
		if r.tasks[i].ID == id {
			r.tasks[i].Status = normalizeTaskStatus(status)
			return nil
		}
	}

	return fmt.Errorf("任务不存在")
}

func (r *memoryStudyRepository) DeleteTask(id uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := range r.tasks {
		if r.tasks[i].ID == id {
			r.tasks = append(r.tasks[:i], r.tasks[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("任务不存在")
}

func (r *memoryStudyRepository) ListRecords(limit int) ([]studyRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	records := append([]studyRecord(nil), r.records...)
	if limit > 0 && len(records) > limit {
		return records[:limit], nil
	}
	return records, nil
}

func (r *memoryStudyRepository) CreateRecord(req createRecordRequest) (studyRecord, studyNote, error) {
	title := strings.TrimSpace(req.Title)
	category := strings.TrimSpace(req.Category)
	summary := strings.TrimSpace(req.Summary)
	noteContent := strings.TrimSpace(req.Note)
	if title == "" || category == "" || summary == "" || noteContent == "" {
		return studyRecord{}, studyNote{}, fmt.Errorf("标题、分类、记录和笔记不能为空")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	record := studyRecord{
		ID:        r.nextID,
		Title:     title,
		Category:  category,
		Summary:   summary,
		CreatedAt: time.Now().Format(timeLayout),
	}
	r.nextID++
	note := studyNote{
		ID:        r.nextID,
		Title:     title + " 笔记",
		Content:   noteContent,
		Tag:       category,
		CreatedAt: time.Now().Format(timeLayout),
	}
	r.nextID++

	r.records = append([]studyRecord{record}, r.records...)
	r.notes = append([]studyNote{note}, r.notes...)

	return record, note, nil
}

func (r *memoryStudyRepository) ListNotes(limit int) ([]studyNote, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	notes := append([]studyNote(nil), r.notes...)
	if limit > 0 && len(notes) > limit {
		return notes[:limit], nil
	}
	return notes, nil
}

func (r *memoryStudyRepository) UpdateRecord(id uint, req updateRecordRequest) (studyRecord, error) {
	title := strings.TrimSpace(req.Title)
	category := strings.TrimSpace(req.Category)
	summary := strings.TrimSpace(req.Summary)
	if title == "" || category == "" || summary == "" {
		return studyRecord{}, fmt.Errorf("标题、分类和记录不能为空")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for i := range r.records {
		if r.records[i].ID == id {
			r.records[i].Title = title
			r.records[i].Category = category
			r.records[i].Summary = summary
			return r.records[i], nil
		}
	}

	return studyRecord{}, fmt.Errorf("记录不存在")
}

func (r *memoryStudyRepository) DeleteRecord(id uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := range r.records {
		if r.records[i].ID == id {
			r.records = append(r.records[:i], r.records[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("记录不存在")
}

func (r *memoryStudyRepository) UpdateNote(id uint, req updateNoteRequest) (studyNote, error) {
	title := strings.TrimSpace(req.Title)
	tag := strings.TrimSpace(req.Tag)
	content := strings.TrimSpace(req.Content)
	if title == "" || tag == "" || content == "" {
		return studyNote{}, fmt.Errorf("标题、标签和笔记不能为空")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for i := range r.notes {
		if r.notes[i].ID == id {
			r.notes[i].Title = title
			r.notes[i].Tag = tag
			r.notes[i].Content = content
			return r.notes[i], nil
		}
	}

	return studyNote{}, fmt.Errorf("笔记不存在")
}

func (r *memoryStudyRepository) DeleteNote(id uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := range r.notes {
		if r.notes[i].ID == id {
			r.notes = append(r.notes[:i], r.notes[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("笔记不存在")
}

const timeLayout = "2006-01-02 15:04"

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	data, err := os.ReadFile("./index.html")
	if err != nil {
		http.Error(w, "内部错误", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/html; charset=utf-8")
	fmt.Fprint(w, string(data))
}

func StudyDashboardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data, err := buildDashboard()
	if err != nil {
		http.Error(w, "内部错误", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, studyResponse{Code: 0, Data: data})
}

func StudyTaskListHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tasks, err := studyRepo.ListTasks()
		if err != nil {
			http.Error(w, "内部错误", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{"code": 0, "data": tasks})
	case http.MethodPost:
		var req createTaskRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "参数错误", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		task, err := studyRepo.CreateTask(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		writeJSON(w, http.StatusCreated, map[string]interface{}{"code": 0, "data": task})
	case http.MethodDelete:
		id, err := parseUintQuery(r, "id")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := studyRepo.DeleteTask(id); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{"code": 0})
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func StudyTaskStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := parseUintQuery(r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req updateTaskStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数错误", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := studyRepo.UpdateTaskStatus(id, req.Status); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"code": 0})
}

func StudyRecordHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		limit := parseLimit(r)
		records, err := studyRepo.ListRecords(limit)
		if err != nil {
			http.Error(w, "内部错误", http.StatusInternalServerError)
			return
		}
		records = filterRecords(records, r.URL.Query().Get("q"), r.URL.Query().Get("date"))
		writeJSON(w, http.StatusOK, map[string]interface{}{"code": 0, "data": records})
	case http.MethodPost:
		var req createRecordRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "参数错误", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		record, note, err := studyRepo.CreateRecord(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		dashboard, err := buildDashboard()
		if err != nil {
			http.Error(w, "内部错误", http.StatusInternalServerError)
			return
		}

		writeJSON(w, http.StatusCreated, map[string]interface{}{
			"code": 0,
			"data": map[string]interface{}{
				"record":    record,
				"note":      note,
				"dashboard": dashboard,
			},
		})
	case http.MethodPut:
		id, err := parseUintQuery(r, "id")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var req updateRecordRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "参数错误", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		record, err := studyRepo.UpdateRecord(id, req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{"code": 0, "data": record})
	case http.MethodDelete:
		id, err := parseUintQuery(r, "id")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := studyRepo.DeleteRecord(id); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{"code": 0})
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func StudyNoteListHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		limit := parseLimit(r)
		notes, err := studyRepo.ListNotes(limit)
		if err != nil {
			http.Error(w, "内部错误", http.StatusInternalServerError)
			return
		}
		notes = filterNotes(notes, r.URL.Query().Get("q"), r.URL.Query().Get("date"))
		writeJSON(w, http.StatusOK, map[string]interface{}{"code": 0, "data": notes})
	case http.MethodPut:
		id, err := parseUintQuery(r, "id")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var req updateNoteRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "参数错误", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		note, err := studyRepo.UpdateNote(id, req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{"code": 0, "data": note})
	case http.MethodDelete:
		id, err := parseUintQuery(r, "id")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := studyRepo.DeleteNote(id); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{"code": 0})
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func buildDashboard() (studyData, error) {
	tasks, err := studyRepo.ListTasks()
	if err != nil {
		return studyData{}, err
	}
	allRecords, err := studyRepo.ListRecords(0)
	if err != nil {
		return studyData{}, err
	}
	records := allRecords
	if len(records) > 5 {
		records = records[:5]
	}
	notes, err := studyRepo.ListNotes(5)
	if err != nil {
		return studyData{}, err
	}

	return studyData{
		App: studyApp{
			Name:    "学习打卡本",
			Tagline: "每日任务、学习记录、复盘笔记集中查看",
		},
		Profile: studyProfile{
			UserName: "今日学习者",
			Streak:   12,
			Total:    len(allRecords),
			Goal:     "连续 30 天保持每天复盘",
		},
		Today: studyToday{
			Date:       time.Now().Format("2006-01-02"),
			Motto:      "先完成，再优化；先记录，再总结。",
			Completion: buildCompletion(tasks),
			Tasks:      tasks,
		},
		RecentRecords: records,
		Notes:         notes,
		WeeklyFocus: []string{
			"每天完成 1 次任务打卡",
			"每天补 1 条学习记录",
			"每晚写 1 条复盘笔记",
		},
		Footer: studyFooter{
			Nav: []string{"今日", "记录", "笔记", "我的"},
		},
	}, nil
}

func buildCompletion(tasks []studyTask) string {
	if len(tasks) == 0 {
		return "0 / 0"
	}

	done := 0
	for _, task := range tasks {
		if task.Status == "已完成" {
			done++
		}
	}

	return fmt.Sprintf("%d / %d", done, len(tasks))
}

func normalizeTaskStatus(status string) string {
	switch strings.TrimSpace(status) {
	case "已完成", "进行中", "待开始":
		return status
	default:
		return "待开始"
	}
}

func filterRecords(records []studyRecord, keyword string, date string) []studyRecord {
	keyword = strings.TrimSpace(strings.ToLower(keyword))
	date = strings.TrimSpace(date)
	filtered := make([]studyRecord, 0, len(records))

	for _, record := range records {
		if keyword != "" {
			target := strings.ToLower(record.Title + " " + record.Category + " " + record.Summary)
			if !strings.Contains(target, keyword) {
				continue
			}
		}
		if date != "" && !strings.HasPrefix(record.CreatedAt, date) {
			continue
		}
		filtered = append(filtered, record)
	}

	return filtered
}

func filterNotes(notes []studyNote, keyword string, date string) []studyNote {
	keyword = strings.TrimSpace(strings.ToLower(keyword))
	date = strings.TrimSpace(date)
	filtered := make([]studyNote, 0, len(notes))

	for _, note := range notes {
		if keyword != "" {
			target := strings.ToLower(note.Title + " " + note.Tag + " " + note.Content)
			if !strings.Contains(target, keyword) {
				continue
			}
		}
		if date != "" && !strings.HasPrefix(note.CreatedAt, date) {
			continue
		}
		filtered = append(filtered, note)
	}

	return filtered
}

func parseLimit(r *http.Request) int {
	value := r.URL.Query().Get("limit")
	if value == "" {
		return 0
	}

	limit, err := strconv.Atoi(value)
	if err != nil || limit < 0 {
		return 0
	}
	return limit
}

func parseUintQuery(r *http.Request, key string) (uint, error) {
	value := strings.TrimSpace(r.URL.Query().Get(key))
	if value == "" {
		return 0, fmt.Errorf("缺少 %s 参数", key)
	}

	num, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%s 参数错误", key)
	}
	return uint(num), nil
}

func mapTasks(items []model.StudyTaskModel) []studyTask {
	tasks := make([]studyTask, 0, len(items))
	for _, item := range items {
		tasks = append(tasks, mapTask(item))
	}
	return tasks
}

func mapTask(item model.StudyTaskModel) studyTask {
	return studyTask{
		ID:       item.ID,
		Title:    item.Title,
		Duration: item.Duration,
		Status:   item.Status,
	}
}

func mapRecords(items []model.StudyRecordModel) []studyRecord {
	records := make([]studyRecord, 0, len(items))
	for _, item := range items {
		records = append(records, mapRecord(item))
	}
	return records
}

func mapRecord(item model.StudyRecordModel) studyRecord {
	return studyRecord{
		ID:        item.ID,
		Title:     item.Title,
		Category:  item.Category,
		Summary:   item.Summary,
		CreatedAt: item.CreatedAt.Format(timeLayout),
	}
}

func mapNotes(items []model.StudyNoteModel) []studyNote {
	notes := make([]studyNote, 0, len(items))
	for _, item := range items {
		notes = append(notes, mapNote(item))
	}
	return notes
}

func mapNote(item model.StudyNoteModel) studyNote {
	return studyNote{
		ID:        item.ID,
		Title:     item.Title,
		Content:   item.Content,
		Tag:       item.Tag,
		CreatedAt: item.CreatedAt.Format(timeLayout),
	}
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	body, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "内部错误", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	fmt.Fprint(w, string(body))
}
