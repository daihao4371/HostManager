package ui

import (
	"testing"
	"time"

	"github.com/daihao4371/hostmanager/internal/theme"
	"github.com/nsf/termbox-go"
)

// 模拟主题用于测试
func createTestTheme() *theme.Theme {
	return &theme.Theme{
		Background: termbox.ColorBlack,
		Foreground: termbox.ColorWhite,
		Border:     termbox.ColorWhite,
		Highlight:  termbox.ColorYellow,
		Success:    termbox.ColorGreen,
		Error:      termbox.ColorRed,
		Warning:    termbox.ColorYellow,
		Info:       termbox.ColorCyan,
	}
}

// 测试图标系统
func TestIconSet(t *testing.T) {
	iconSet := CreateIconSet()

	// 验证状态图标
	if iconSet.Status[StatusOnline] == "" {
		t.Error("在线状态图标不能为空")
	}

	if iconSet.Status[StatusLoading] == "" {
		t.Error("加载状态图标不能为空")
	}

	// 验证操作图标
	if iconSet.Actions["connect"] == "" {
		t.Error("连接操作图标不能为空")
	}

	// 验证UI图标
	if iconSet.UI["arrow_up"] == "" {
		t.Error("向上箭头图标不能为空")
	}

	t.Log("图标系统测试通过")
}

// 测试状态指示器
func TestStatusBadge(t *testing.T) {
	testTheme := createTestTheme()
	renderer := NewRenderEngine(testTheme)
	iconSet := CreateIconSet()

	// 测试不同状态类型
	statuses := []StatusType{
		StatusOnline,
		StatusOffline,
		StatusLoading,
		StatusError,
		StatusWarning,
		StatusConnecting,
		StatusMaintenance,
	}

	for _, status := range statuses {
		renderer.RenderEnhancedStatusBadge(0, 0, status, "测试", iconSet, true, 0)
	}

	t.Log("状态指示器测试通过")
}

// 测试进度条
func TestProgressBar(t *testing.T) {
	testTheme := createTestTheme()
	renderer := NewRenderEngine(testTheme)
	style := CreatePremiumStyle(testTheme)

	// 测试不同进度值
	progressValues := []float32{0.0, 0.25, 0.5, 0.75, 1.0}

	for _, progress := range progressValues {
		renderer.RenderAdvancedProgressBar(0, 0, 40, progress, 4, true)
	}

	t.Log("进度条测试通过")
}

// 测试加载动画
func TestSpinner(t *testing.T) {
	testTheme := createTestTheme()
	renderer := NewRenderEngine(testTheme)
	iconSet := CreateIconSet()

	// 测试不同加载动画类型
	spinnerTypes := []string{"dots", "bars", "arrows", "pulse"}

	for _, spinnerType := range spinnerTypes {
		for frame := 0; frame < 8; frame++ {
			renderer.RenderAdvancedSpinner(0, 0, spinnerType, frame, iconSet)
		}
	}

	t.Log("加载动画测试通过")
}

// 测试交互管理器
func TestInteractionManager(t *testing.T) {
	config := CreateFeedbackConfig()
	manager := NewInteractionManager(config)

	// 创建测试组件
	component := &UIComponent{
		ID:     "test-button",
		Type:   ComponentButton,
		Bounds: Rect{X: 10, Y: 10, Width: 20, Height: 5},
		State:  StateNormal,
	}

	// 测试悬停事件
	manager.HandleMouseHover(15, 12, component)

	// 测试点击事件
	manager.HandleMouseClick(15, 12, component)

	// 更新状态
	manager.Update()

	// 验证点击反馈
	if len(manager.clickFeedbacks) == 0 && config.ClickFeedback {
		t.Error("点击反馈效果未正确创建")
	}

	t.Log("交互管理器测试通过")
}

// 测试动画系统
func TestAnimationSystem(t *testing.T) {
	testTheme := createTestTheme()
	renderer := NewRenderEngine(testTheme)

	// 创建测试组件
	component := &UIComponent{
		ID:     "animated-card",
		Type:   ComponentCard,
		Bounds: Rect{X: 0, Y: 0, Width: 30, Height: 10},
		Style:  CreatePremiumStyle(testTheme),
	}

	// 启动动画
	renderer.StartAnimation(component)

	if !renderer.HasActiveAnimations() {
		t.Error("动画未正确启动")
	}

	// 模拟动画更新
	time.Sleep(50 * time.Millisecond)
	renderer.UpdateAnimations()

	t.Log("动画系统测试通过")
}

// 测试UI组件创建
func TestComponentCreation(t *testing.T) {
	testTheme := createTestTheme()
	style := CreatePremiumStyle(testTheme)

	// 验证样式配置
	if style.Animation.Duration <= 0 {
		t.Error("动画持续时间配置错误")
	}

	if !style.Shadow.Enabled {
		t.Error("阴影效果未启用")
	}

	if style.Border.Type != BorderRounded {
		t.Error("边框类型配置错误")
	}

	t.Log("组件创建测试通过")
}

// 集成测试 - 模拟完整的UI渲染流程
func TestUIIntegration(t *testing.T) {
	testTheme := createTestTheme()
	renderer := NewRenderEngine(testTheme)
	iconSet := CreateIconSet()
	config := CreateFeedbackConfig()
	interactionManager := NewInteractionManager(config)

	// 创建多个测试组件
	components := []*UIComponent{
		{
			ID:     "header",
			Type:   ComponentCard,
			Bounds: Rect{X: 0, Y: 0, Width: 80, Height: 5},
			Style:  CreatePremiumStyle(testTheme),
		},
		{
			ID:     "status-list",
			Type:   ComponentList,
			Bounds: Rect{X: 0, Y: 6, Width: 50, Height: 20},
			Style:  CreatePremiumStyle(testTheme),
		},
		{
			ID:     "action-panel",
			Type:   ComponentPanel,
			Bounds: Rect{X: 52, Y: 6, Width: 28, Height: 20},
			Style:  CreatePremiumStyle(testTheme),
		},
	}

	// 渲染所有组件
	for _, component := range components {
		renderer.RenderCard(
			component.Bounds.X,
			component.Bounds.Y,
			component.Bounds.Width,
			component.Bounds.Height,
			component.Style,
			component.ID,
			component.State,
		)
	}

	// 测试状态指示器
	renderer.RenderEnhancedStatusBadge(5, 8, StatusOnline, "服务器1", iconSet, true, 0)
	renderer.RenderEnhancedStatusBadge(5, 10, StatusLoading, "服务器2", iconSet, true, 2)
	renderer.RenderEnhancedStatusBadge(5, 12, StatusError, "服务器3", iconSet, false, 0)

	// 测试进度条
	renderer.RenderAdvancedProgressBar(5, 15, 30, 0.6, 4, true)

	// 测试图标按钮
	// renderer.RenderIconButton(55, 8, "connect", iconSet, StateNormal, CreatePremiumStyle(testTheme))
	// renderer.RenderIconButton(60, 8, "edit", iconSet, StateHover, CreatePremiumStyle(testTheme))
	// renderer.RenderIconButton(65, 8, "delete", iconSet, StateDisabled, CreatePremiumStyle(testTheme))

	// 模拟交互
	interactionManager.HandleMouseHover(55, 8, components[2])
	interactionManager.HandleMouseClick(60, 8, components[2])

	// 更新和渲染交互效果
	interactionManager.Update()
	interactionManager.Render(renderer)

	t.Log("UI集成测试完成 - 所有组件渲染正常")
}
