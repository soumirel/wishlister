package presenter

import (
	"errors"
	"fmt"

	"github.com/soumirel/wishlister/services/telegram-bot/internal/domain/ui"
)

type presenterProvider struct {
	presenters map[ui.ModuleType]ui.Presenter
}

func NewPresenterProvider(
	presenters ...ui.Presenter,
) (*presenterProvider, error) {
	p := &presenterProvider{
		presenters: make(map[ui.ModuleType]ui.Presenter),
	}
	err := p.initPresenters(presenters...)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (p *presenterProvider) initPresenters(presenters ...ui.Presenter) error {
	for _, pr := range presenters {
		module := pr.Module()
		if _, ok := p.presenters[module]; ok {
			return errors.New("met presenter module duplicate")
		}
		p.presenters[pr.Module()] = pr
	}
	return nil
}

func (p *presenterProvider) GetPresenter(m ui.ModuleType) (ui.Presenter, error) {
	presenter, ok := p.presenters[m]
	if !ok {
		return nil, fmt.Errorf("no registered presenters for module")
	}
	return presenter, nil
}
