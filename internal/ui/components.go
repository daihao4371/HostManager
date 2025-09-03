package ui

import (
	"time"

	"github.com/daihao4371/hostmanager/internal/theme"
	"github.com/nsf/termbox-go"
)

// 动画状态枚举
type AnimationType int

const (
	AnimationNone       AnimationType = iota
	AnimationFade                     // 淡入淡出
	AnimationSlide                    // 滑动
	AnimationPulse                    // 脉冲
	AnimationTypewriter               // 打字机效果
)

// UI组件基础样式
type ComponentStyle struct {
	Padding    Padding         // 内边距
	Margin     Margin          // 外边距
	Border     BorderStyle     // 边框样式
	Background BackgroundStyle // 背景样式
	Shadow     ShadowStyle     // 阴影效果
	Animation  AnimationConfig // 动画配置
}

// 内边距定义
type Padding struct {
	Top, Right, Bottom, Left int
}

// 外边距定义
type Margin struct {
	Top, Right, Bottom, Left int
}

// 边框样式
type BorderStyle struct {
	Type   BorderType        // 边框类型
	Width  int               // 边框宽度
	Color  termbox.Attribute // 边框颜色
	Radius int               // 圆角半径（模拟）
}

// 边框类型枚举
type BorderType int

const (
	BorderNone    BorderType = iota
	BorderSolid              // 实线
	BorderDashed             // 虚线
	BorderDouble             // 双线
	BorderRounded            // 圆角
)

// 背景样式
type BackgroundStyle struct {
	Type        BackgroundType    // 背景类型
	Color       termbox.Attribute // 主色
	GradientEnd termbox.Attribute // 渐变终点色
	Pattern     string            // 背景图案
	Opacity     float32           // 透明度
}

// 背景类型枚举
type BackgroundType int

const (
	BackgroundSolid    BackgroundType = iota
	BackgroundGradient                // 渐变
	BackgroundPattern                 // 图案
)

// 阴影样式
type ShadowStyle struct {
	Enabled bool              // 是否启用
	OffsetX int               // X偏移
	OffsetY int               // Y偏移
	Blur    int               // 模糊半径
	Color   termbox.Attribute // 阴影颜色
	Opacity float32           // 透明度
}

// 动画配置
type AnimationConfig struct {
	Type     AnimationType // 动画类型
	Duration time.Duration // 持续时间
	Delay    time.Duration // 延迟时间
	Easing   EasingType    // 缓动函数
	Loop     bool          // 是否循环
}

// 缓动函数类型
type EasingType int

const (
	EasingLinear EasingType = iota
	EasingEaseIn
	EasingEaseOut
	EasingEaseInOut
	EasingBounce
)

// UI组件基类
type UIComponent struct {
	ID        string         // 组件ID
	Type      ComponentType  // 组件类型
	Bounds    Rect           // 组件边界
	Style     ComponentStyle // 样式配置
	State     ComponentState // 组件状态
	Visible   bool           // 是否可见
	Enabled   bool           // 是否启用
	Focused   bool           // 是否聚焦
	Animation *Animation     // 当前动画
}

// 组件类型枚举
type ComponentType int

const (
	ComponentCard ComponentType = iota
	ComponentButton
	ComponentList
	ComponentPanel
	ComponentModal
	ComponentToast
	ComponentBadge
	ComponentSpinner
)

// 组件状态枚举
type ComponentState int

const (
	StateNormal ComponentState = iota
	StateHover
	StatePressed
	StateFocused
	StateDisabled
	StateLoading
)

// 矩形定义
type Rect struct {
	X, Y, Width, Height int
}

// 动画实例
type Animation struct {
	Config    AnimationConfig
	StartTime time.Time
	Progress  float32 // 0.0 - 1.0
	Active    bool
}

// 状态类型扩展
type StatusType int

const (
	StatusOnline StatusType = iota
	StatusOffline
	StatusLoading
	StatusError
	StatusWarning
	StatusIdle
	StatusConnecting
	StatusMaintenance
)

// 图标映射系统
type IconSet struct {
	Status  map[StatusType]string
	Actions map[string]string
	UI      map[string]string
}

// 创建默认图标集
func CreateIconSet() *IconSet {
	return &IconSet{
		Status: map[StatusType]string{
			StatusOnline:      "●",
			StatusOffline:     "○",
			StatusLoading:     "◐",
			StatusError:       "✗",
			StatusWarning:     "⚠",
			StatusIdle:        "◌",
			StatusConnecting:  "◔",
			StatusMaintenance: "🔧",
		},
		Actions: map[string]string{
			"connect":    "🔗",
			"disconnect": "⛔",
			"edit":       "✏",
			"delete":     "🗑",
			"add":        "➕",
			"refresh":    "🔄",
			"save":       "💾",
			"settings":   "⚙",
		},
		UI: map[string]string{
			"arrow_up":    "▲",
			"arrow_down":  "▼",
			"arrow_left":  "◀",
			"arrow_right": "▶",
			"check":       "✓",
			"cross":       "✗",
			"menu":        "☰",
			"search":      "🔍",
		},
	}
}

// 交互反馈配置
type FeedbackConfig struct {
	HoverDelay    time.Duration // 悬停延迟
	ClickFeedback bool          // 点击反馈
	SoundEnabled  bool          // 声音反馈
	HapticEnabled bool          // 触觉反馈
}

// 创建默认交互反馈配置
func CreateFeedbackConfig() FeedbackConfig {
	return FeedbackConfig{
		HoverDelay:    100 * time.Millisecond,
		ClickFeedback: true,
		SoundEnabled:  false,
		HapticEnabled: false,
	}
}

// 创建默认的高级样式
func CreatePremiumStyle(theme *theme.Theme) ComponentStyle {
	return ComponentStyle{
		Padding: Padding{Top: 1, Right: 2, Bottom: 1, Left: 2},
		Margin:  Margin{Top: 0, Right: 0, Bottom: 1, Left: 0},
		Border: BorderStyle{
			Type:   BorderRounded,
			Width:  1,
			Color:  theme.Border,
			Radius: 2,
		},
		Background: BackgroundStyle{
			Type:        BackgroundGradient,
			Color:       theme.Background,
			GradientEnd: theme.Info,
			Opacity:     0.95,
		},
		Shadow: ShadowStyle{
			Enabled: true,
			OffsetX: 1,
			OffsetY: 1,
			Blur:    2,
			Color:   termbox.ColorBlack,
			Opacity: 0.3,
		},
		Animation: AnimationConfig{
			Type:     AnimationFade,
			Duration: 200 * time.Millisecond,
			Easing:   EasingEaseOut,
		},
	}
}

// 应用缓动函数
func applyEasing(progress float32, easing EasingType) float32 {
	switch easing {
	case EasingEaseIn:
		return progress * progress
	case EasingEaseOut:
		return 1 - (1-progress)*(1-progress)
	case EasingEaseInOut:
		if progress < 0.5 {
			return 2 * progress * progress
		}
		return 1 - 2*(1-progress)*(1-progress)
	case EasingBounce:
		if progress < 1.0/2.75 {
			return 7.5625 * progress * progress
		} else if progress < 2.0/2.75 {
			progress -= 1.5 / 2.75
			return 7.5625*progress*progress + 0.75
		}
		progress -= 2.25 / 2.75
		return 7.5625*progress*progress + 0.9375
	default:
		return progress
	}
}
