const {
  MAIN_TYPES,
  STYLE_TYPES,
  MAIN_QUESTIONS,
  STYLE_QUESTIONS,
  buildResultCopy,
} = require('../../data/quiz')

const PROGRESS_COPY = [
  '你已经击败 37% 的嘴硬同事，继续暴露本性。',
  '别慌，测出来不是有病，只是有工龄。',
  '目前看你还算克制，再答几题就原形毕露了。',
  '你的工位灵魂正在加载，请继续提交证词。',
  '再坚持一下，职场宇宙快把你归类了。',
]

const STORAGE_KEY = 'workplace-persona-last-result'

Page({
  data: {
    stage: 'landing',
    introTags: ['36 题主测', '9 题风格加试', '108 种职场角色'],
    highlights: [
      '题目专治会议创伤、群聊阴阳怪气和周一早会后遗症。',
      '结果不是正经心理学诊断，是打工人自我认领现场。',
      '测完直接给你主人格、副人格和专属工位称号。',
    ],
    roleHighlights: buildRoleCloud(),
    mainTypes: MAIN_TYPES,
    styleTypes: STYLE_TYPES,
    totalMain: MAIN_QUESTIONS.length,
    totalStyle: STYLE_QUESTIONS.length,
    totalQuestions: MAIN_QUESTIONS.length + STYLE_QUESTIONS.length,
    questionQueue: [],
    currentQuestion: null,
    progressPercent: 0,
    currentIndex: 0,
    currentPhaseLabel: '主测试',
    answerLabel: '',
    selectedOptionKey: '',
    scoreBoard: createEmptyBoard(MAIN_TYPES),
    styleBoard: createEmptyBoard(STYLE_TYPES),
    answers: [],
    result: null,
    lastResult: null,
    posterImage: '',
    posterVisible: false,
    generatingPoster: false,
  },

  onLoad() {
    this.loadLastResult()
  },

  onShareAppMessage() {
    const { result } = this.data
    if (!result) {
      return {
        title: '工位人格研究所：测测你是哪种职场角色',
        path: '/pages/index/index',
      }
    }

    return {
      title: `我的工位人格是 ${result.headline}，你也来测测`,
      path: '/pages/index/index',
    }
  },

  onShareTimeline() {
    const { result } = this.data
    return {
      title: result ? `我的工位人格是 ${result.headline}` : '测测你的工位人格',
      query: '',
    }
  },

  startQuiz() {
    this.resetState('quiz')
    this.setData({
      questionQueue: buildQuestionQueue(),
    })
    this.loadQuestion(0)
  },

  restartQuiz() {
    this.resetState('landing')
    this.loadLastResult()
  },

  chooseOption(event) {
    const { key } = event.currentTarget.dataset
    const { currentQuestion } = this.data

    if (!currentQuestion) {
      return
    }

    const option = currentQuestion.options.find((item) => item.key === key)
    if (!option) {
      return
    }

    const nextAnswers = [...this.data.answers, {
      questionId: currentQuestion.id,
      questionTitle: currentQuestion.title,
      phase: currentQuestion.phase,
      optionKey: option.key,
      optionText: option.text,
      scoreKey: option.scoreKey,
    }]

    if (currentQuestion.phase === 'main') {
      this.data.scoreBoard[option.scoreKey] += 1
    } else {
      this.data.styleBoard[option.scoreKey] += 1
    }

    this.setData({
      answers: nextAnswers,
      selectedOptionKey: option.key,
      answerLabel: option.text,
      scoreBoard: this.data.scoreBoard,
      styleBoard: this.data.styleBoard,
    })

    const nextIndex = this.data.currentIndex + 1
    if (nextIndex < this.data.totalQuestions) {
      this.loadQuestion(nextIndex)
      return
    }

    this.finishQuiz()
  },

  loadQuestion(index) {
    const allQuestions = this.data.questionQueue.length ? this.data.questionQueue : buildQuestionQueue()
    const currentQuestion = allQuestions[index]
    const progressPercent = Math.round((index / allQuestions.length) * 100)
    const currentPhaseLabel = currentQuestion.phase === 'main' ? '主测试' : '风格加试'
    const answerLabel = PROGRESS_COPY[index % PROGRESS_COPY.length]

    this.setData({
      stage: 'quiz',
      questionQueue: allQuestions,
      currentQuestion,
      currentIndex: index,
      progressPercent,
      currentPhaseLabel,
      answerLabel,
      selectedOptionKey: '',
    })
  },

  finishQuiz() {
    const rankedMain = rankBoard(this.data.scoreBoard, MAIN_TYPES)
    const rankedStyle = rankBoard(this.data.styleBoard, STYLE_TYPES)
    const primary = rankedMain[0]
    const secondary = rankedMain[1] || rankedMain[0]
    const style = rankedStyle[0]
    const result = buildResultCopy(primary.key, style.key, secondary.key)
    const finalResult = {
      ...result,
      primary,
      secondary,
      style,
      rankedMain,
      rankedStyle,
      topStrengthTypes: rankedMain.slice(0, 3),
      shareLine: `我是 ${result.headline}，表面 ${primary.label}，骨子里还藏着 ${secondary.label}。`,
      posterFacts: [
        `${primary.label} 主人格`,
        `${style.label} 出招风格`,
        `${secondary.label} 隐藏副人格`,
      ],
    }

    this.setData({
      stage: 'result',
      progressPercent: 100,
      result: finalResult,
      lastResult: finalResult,
    })
    wx.setStorageSync(STORAGE_KEY, finalResult)
  },

  copySummary() {
    const { result } = this.data
    if (!result) {
      return
    }

    const text = [
      `我的职场人格是：${result.headline}`,
      `${result.subtitle}`,
      result.verdict,
      `副人格：${result.secondary.label}`,
      result.roast,
    ].join('\n')

    wx.setClipboardData({
      data: text,
      success: () => {
        wx.showToast({ title: '结果文案已复制', icon: 'success' })
      },
    })
  },

  async generatePoster() {
    const { result, generatingPoster } = this.data
    if (!result || generatingPoster) {
      return
    }

    this.setData({ generatingPoster: true })

    try {
      const imagePath = await drawPoster(this, result)
      this.setData({
        posterImage: imagePath,
        posterVisible: true,
      })
    } catch (error) {
      wx.showToast({ title: '海报生成失败', icon: 'none' })
    } finally {
      this.setData({ generatingPoster: false })
    }
  },

  previewPoster() {
    const { posterImage } = this.data
    if (!posterImage) {
      return
    }

    wx.previewImage({
      urls: [posterImage],
      current: posterImage,
    })
  },

  savePoster() {
    const { posterImage } = this.data
    if (!posterImage) {
      wx.showToast({ title: '请先生成海报', icon: 'none' })
      return
    }

    wx.saveImageToPhotosAlbum({
      filePath: posterImage,
      success: () => {
        wx.showToast({ title: '已保存到相册', icon: 'success' })
      },
      fail: () => {
        wx.showToast({ title: '保存失败，请检查相册权限', icon: 'none' })
      },
    })
  },

  closePoster() {
    this.setData({ posterVisible: false })
  },

  resetState(stage) {
    this.setData({
      stage,
      questionQueue: [],
      currentQuestion: null,
      progressPercent: 0,
      currentIndex: 0,
      currentPhaseLabel: '主测试',
      answerLabel: '',
      selectedOptionKey: '',
      scoreBoard: createEmptyBoard(MAIN_TYPES),
      styleBoard: createEmptyBoard(STYLE_TYPES),
      answers: [],
      result: null,
      posterImage: '',
      posterVisible: false,
      generatingPoster: false,
    })
  },

  loadLastResult() {
    const lastResult = wx.getStorageSync(STORAGE_KEY)
    if (!lastResult) {
      return
    }

    this.setData({ lastResult })
  },
})

function getAllQuestions() {
  return [
    ...MAIN_QUESTIONS.map((item) => ({ ...item, phase: 'main' })),
    ...STYLE_QUESTIONS.map((item) => ({ ...item, phase: 'style' })),
  ]
}

function buildQuestionQueue() {
  const mainQuestions = shuffle(
    MAIN_QUESTIONS.map((item) => ({
      ...item,
      phase: 'main',
      options: shuffle(item.options),
    })),
  )
  const styleQuestions = shuffle(
    STYLE_QUESTIONS.map((item) => ({
      ...item,
      phase: 'style',
      options: shuffle(item.options),
    })),
  )

  return [...mainQuestions, ...styleQuestions]
}

function rankBoard(board, dictionary) {
  return dictionary
    .map((item) => ({
      key: item.key,
      label: item.label,
      score: board[item.key] || 0,
      summary: item.summary || item.tone,
    }))
    .sort((a, b) => b.score - a.score)
}

function createEmptyBoard(dictionary) {
  return dictionary.reduce((acc, item) => {
    acc[item.key] = 0
    return acc
  }, {})
}

function buildRoleCloud() {
  return [
    '破晓冲锋手',
    '分歧拆弹员',
    '趋势瞭望塔',
    '脑洞喷焰兽',
    '贵人雷达站',
    '定海压舱石',
  ]
}

function shuffle(list) {
  const cloned = [...list]
  for (let index = cloned.length - 1; index > 0; index -= 1) {
    const randomIndex = Math.floor(Math.random() * (index + 1))
    const current = cloned[index]
    cloned[index] = cloned[randomIndex]
    cloned[randomIndex] = current
  }
  return cloned
}

async function drawPoster(page, result) {
  const query = page.createSelectorQuery()
  const canvasNode = await new Promise((resolve, reject) => {
    query.select('#posterCanvas').fields({ node: true, size: true }).exec((res) => {
      const target = res && res[0]
      if (!target || !target.node) {
        reject(new Error('canvas node missing'))
        return
      }
      resolve(target)
    })
  })

  const { node: canvas, width, height } = canvasNode
  const ctx = canvas.getContext('2d')
  const systemInfo = wx.getWindowInfo ? wx.getWindowInfo() : wx.getSystemInfoSync()
  const dpr = systemInfo.pixelRatio || 2

  canvas.width = width * dpr
  canvas.height = height * dpr
  ctx.scale(dpr, dpr)

  drawPosterBackground(ctx, width, height)
  drawPosterCard(ctx, width, height)

  ctx.fillStyle = '#FFE5D5'
  ctx.font = '12px sans-serif'
  ctx.fillText('WORKPLACE PERSONA POSTER', 26, 38)

  ctx.fillStyle = '#FFFFFF'
  ctx.font = 'bold 34px sans-serif'
  fillTextLines(ctx, result.headline, 26, 84, 48, width - 52, 2)

  ctx.fillStyle = 'rgba(255, 244, 236, 0.78)'
  ctx.font = '16px sans-serif'
  fillTextLines(ctx, result.subtitle, 26, 152, 26, width - 52, 1)

  drawPills(ctx, result.tags, 26, 188, width - 52)

  ctx.fillStyle = 'rgba(255, 251, 247, 0.18)'
  roundRect(ctx, 26, 246, width - 52, 124, 22)
  ctx.fill()

  ctx.fillStyle = '#FFF7F1'
  ctx.font = '16px sans-serif'
  fillTextLines(ctx, result.shareLine, 42, 280, 28, width - 84, 3)

  ctx.fillStyle = '#2A2452'
  roundRect(ctx, 26, 390, width - 52, 142, 24)
  ctx.fill()

  ctx.fillStyle = '#FFC49D'
  ctx.font = '12px sans-serif'
  ctx.fillText('人格标签', 42, 420)

  ctx.fillStyle = '#FFFFFF'
  ctx.font = 'bold 22px sans-serif'
  ctx.fillText(`${result.primary.label} / ${result.style.label}`, 42, 454)

  ctx.fillStyle = 'rgba(255,255,255,0.82)'
  ctx.font = '14px sans-serif'
  fillTextLines(ctx, result.verdict, 42, 486, 24, width - 84, 2)

  ctx.fillStyle = 'rgba(42, 36, 82, 0.88)'
  roundRect(ctx, 26, 552, width - 52, 146, 24)
  ctx.fill()

  ctx.fillStyle = '#FFD5BC'
  ctx.font = '12px sans-serif'
  ctx.fillText('高光技能', 42, 582)

  ctx.fillStyle = '#FFFFFF'
  ctx.font = '15px sans-serif'
  result.strengths.slice(0, 2).forEach((item, index) => {
    fillBulletLine(ctx, item, 42, 614 + index * 34, width - 84)
  })

  ctx.fillStyle = 'rgba(255, 246, 238, 0.72)'
  ctx.font = '12px sans-serif'
  ctx.fillText('工位人格研究所 · 发给同事看看谁先认领', 26, height - 26)

  return new Promise((resolve, reject) => {
    wx.canvasToTempFilePath({
      canvas,
      fileType: 'png',
      success: (res) => resolve(res.tempFilePath),
      fail: reject,
    })
  })
}

function drawPosterBackground(ctx, width, height) {
  const gradient = ctx.createLinearGradient(0, 0, width, height)
  gradient.addColorStop(0, '#1E2042')
  gradient.addColorStop(0.55, '#3D2A69')
  gradient.addColorStop(1, '#FF8357')
  ctx.fillStyle = gradient
  ctx.fillRect(0, 0, width, height)

  const glowA = ctx.createRadialGradient(width - 40, 60, 10, width - 40, 60, 180)
  glowA.addColorStop(0, 'rgba(255, 202, 146, 0.34)')
  glowA.addColorStop(1, 'rgba(255, 202, 146, 0)')
  ctx.fillStyle = glowA
  ctx.fillRect(0, 0, width, height)

  const glowB = ctx.createRadialGradient(40, height - 120, 10, 40, height - 120, 140)
  glowB.addColorStop(0, 'rgba(104, 146, 255, 0.22)')
  glowB.addColorStop(1, 'rgba(104, 146, 255, 0)')
  ctx.fillStyle = glowB
  ctx.fillRect(0, 0, width, height)
}

function drawPosterCard(ctx, width, height) {
  ctx.fillStyle = 'rgba(255,255,255,0.06)'
  roundRect(ctx, 14, 14, width - 28, height - 28, 30)
  ctx.fill()
}

function drawPills(ctx, items, startX, startY, maxWidth) {
  let x = startX
  let y = startY
  const gap = 10
  ctx.font = '13px sans-serif'

  items.forEach((item) => {
    const textWidth = ctx.measureText(item).width
    const pillWidth = textWidth + 26
    if (x + pillWidth > startX + maxWidth) {
      x = startX
      y += 36
    }

    ctx.fillStyle = 'rgba(255,255,255,0.14)'
    roundRect(ctx, x, y, pillWidth, 28, 14)
    ctx.fill()

    ctx.fillStyle = '#FFF4EB'
    ctx.fillText(item, x + 13, y + 19)
    x += pillWidth + gap
  })
}

function fillTextLines(ctx, text, x, y, lineHeight, maxWidth, maxLines) {
  const lines = wrapText(ctx, text, maxWidth, maxLines)
  lines.forEach((line, index) => {
    ctx.fillText(line, x, y + index * lineHeight)
  })
}

function fillBulletLine(ctx, text, x, y, maxWidth) {
  ctx.beginPath()
  ctx.fillStyle = '#FFB26B'
  ctx.arc(x + 4, y - 4, 3, 0, Math.PI * 2)
  ctx.fill()
  ctx.fillStyle = '#FFFFFF'
  ctx.font = '15px sans-serif'
  fillTextLines(ctx, text, x + 14, y, 22, maxWidth - 14, 1)
}

function wrapText(ctx, text, maxWidth, maxLines) {
  const lines = []
  let current = ''

  for (let index = 0; index < text.length; index += 1) {
    const next = current + text[index]
    if (ctx.measureText(next).width > maxWidth) {
      lines.push(current)
      current = text[index]
      if (lines.length === maxLines - 1) {
        break
      }
    } else {
      current = next
    }
  }

  if (lines.length < maxLines && current) {
    lines.push(current)
  }

  if (lines.length === maxLines && text.length > lines.join('').length) {
    const lastIndex = lines.length - 1
    let lastLine = lines[lastIndex]
    while (ctx.measureText(`${lastLine}...`).width > maxWidth && lastLine.length > 0) {
      lastLine = lastLine.slice(0, -1)
    }
    lines[lastIndex] = `${lastLine}...`
  }

  return lines
}

function roundRect(ctx, x, y, width, height, radius) {
  ctx.beginPath()
  ctx.moveTo(x + radius, y)
  ctx.lineTo(x + width - radius, y)
  ctx.quadraticCurveTo(x + width, y, x + width, y + radius)
  ctx.lineTo(x + width, y + height - radius)
  ctx.quadraticCurveTo(x + width, y + height, x + width - radius, y + height)
  ctx.lineTo(x + radius, y + height)
  ctx.quadraticCurveTo(x, y + height, x, y + height - radius)
  ctx.lineTo(x, y + radius)
  ctx.quadraticCurveTo(x, y, x + radius, y)
  ctx.closePath()
}
