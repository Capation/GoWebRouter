package session

import (
	web "Go_Web"
	"github.com/google/uuid"
)

type Manager struct {
	Propagator
	Store
}

func (m *Manager) GetSession(ctx *web.Context) (Session, error) {
	sessID, err := m.Extract(ctx.Req)
	if err != nil {
		return nil, err
	}
	return m.Get(ctx.Req.Context(), sessID)
}

func (m *Manager) InitSession(ctx *web.Context) (Session, error) {
	id := uuid.New().String()
	sess, err := m.Generate(ctx.Req.Context(), id)
	if err != nil {
		return nil, err
	}
	// 注入进去 http响应里面
	err = m.Inject(id, ctx.Resp)
	return sess, err
}

func (m *Manager) RefreshSession(ctx *web.Context) error {
	sess, err := m.GetSession(ctx)
	if err != nil {
		return err
	}
	return m.Refresh(ctx.Req.Context(), sess.ID())
}

func (m *Manager) RemoveSession(ctx *web.Context) error {
	sess, err := m.GetSession(ctx)
	if err != nil {
		return err
	}
	err = m.Store.Remove(ctx.Req.Context(), sess.ID())
	if err != nil {
		return err
	}
	return m.Propagator.Remove(ctx.Resp)
}
