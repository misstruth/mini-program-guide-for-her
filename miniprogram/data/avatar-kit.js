const ARCHETYPE_PRESETS = {
  推进: {
    codename: '油门挂到底的前排发车手',
    fashion: '机能通勤夹克 + 运动感工牌挂绳',
    prop: '荧光排期贴纸和折叠耳机',
    palette: ['#FF845A', '#1C1A3A'],
  },
  统筹: {
    codename: '把混乱排成路线图的控场总调',
    fashion: '利落长风衣 + 模块化腰包',
    prop: '分层任务板和荧光马克笔',
    palette: ['#F6B24F', '#2A2559'],
  },
  执行: {
    codename: '安静但超能打的交付地基',
    fashion: '极简工装衬衫 + 深色机能长裤',
    prop: '金属保温杯和进度清单夹',
    palette: ['#4D7CFF', '#191A35'],
  },
  洞察: {
    codename: '一眼看穿 bug 灵魂的真相显影仪',
    fashion: '高领打底 + 解构西装外套',
    prop: '透明批注板和放大镜吊坠',
    palette: ['#8560FF', '#221B48'],
  },
  创意: {
    codename: '让工位空气突然有梗的灵感暴击怪',
    fashion: '撞色针织 + 手绘贴章外套',
    prop: '涂鸦便签和异形马克杯',
    palette: ['#FF6FAE', '#2C1D49'],
  },
  破局: {
    codename: '专拆死局和旧规矩的带电改造者',
    fashion: '短款皮夹克 + 工装靴',
    prop: '断裂箭头徽章和铆钉文件夹',
    palette: ['#FF7043', '#251A32'],
  },
  协调: {
    codename: '能把吵架群聊重新调成可沟通模式的人',
    fashion: '柔和套头衫 + 叠戴项链',
    prop: '表情贴纸和多色会议卡',
    palette: ['#47C6A3', '#1E2941'],
  },
  稳定: {
    codename: '高压局里永远先把底盘扶住的稳频器',
    fashion: '挺括衬衫 + 廓形马甲',
    prop: '应急预案卡和金属按压笔',
    palette: ['#22B07D', '#1C2842'],
  },
  战略: {
    codename: '总在看半年后地图的方向预言家',
    fashion: '长款大衣 + 未来感窄框眼镜',
    prop: '折叠路线图和星图胸针',
    palette: ['#5A8BFF', '#1E1B52'],
  },
  成长: {
    codename: '把上班过成升级副本的经验狂热者',
    fashion: '学院机能卫衣 + 轻量背包',
    prop: '技能树贴纸和学习打卡章',
    palette: ['#7E6BFF', '#1D1F46'],
  },
  资源: {
    codename: '关键时刻总能摇到人的机会磁场体',
    fashion: '利落西装马甲 + 夸张耳饰',
    prop: '联系人卡册和连接线图钉',
    palette: ['#FF9B54', '#25304E'],
  },
  影响: {
    codename: '能把复杂事讲成全场点头的人形扩音器',
    fashion: '廓形西装 + 高饱和点睛配饰',
    prop: '手持麦克风徽章和提案卡片',
    palette: ['#FF5E7A', '#231C4D'],
  },
}

function buildAvatarBrief(mainKey, styleKey, headline) {
  const preset = ARCHETYPE_PRESETS[mainKey] || ARCHETYPE_PRESETS.推进

  return {
    headline,
    mainKey,
    styleKey,
    codename: preset.codename,
    fashion: preset.fashion,
    prop: preset.prop,
    palette: preset.palette,
  }
}

module.exports = {
  ARCHETYPE_PRESETS,
  buildAvatarBrief,
}
