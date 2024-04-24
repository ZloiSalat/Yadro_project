package pkg

import (
	"YadroProject/types"
	"encoding/json"
	"fmt"
	"net/http"
)

type APIClient struct {
	store Storage
	stem  Stemmer
	cfg   Config
	cmd   *Parser
}

func NewAPIClient(store Storage, stem Stemmer, cfg Config, cmd *Parser) (*APIClient, error) {
	return &APIClient{
		store: store,
		stem:  stem,
		cfg:   cfg,
		cmd:   cmd,
	}, nil
}

func (a *APIClient) GetComicsByID(id int) (types.Comics, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%d/info.0.json", a.cfg.Xkcd.Source, id))
	if err != nil {
		return types.Comics{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return types.Comics{}, fmt.Errorf("status code not OK: %d", resp.StatusCode)
	}

	newData := new(types.Comics)

	if err := json.NewDecoder(resp.Body).Decode(&newData); err != nil {
		return types.Comics{}, err
	}

	return *newData, nil
}

func (a *APIClient) Run() error {

	c, _ := a.cmd.ParseFlag()
	if c {
		fmt.Printf("Config: DBSize %v, Source %v, DBFile %v, End_comics %v, Parallel %v\n", a.cfg.Xkcd.DbSize, a.cfg.Xkcd.Source, a.cfg.Xkcd.DbFile, a.cfg.Xkcd.EndComics, a.cfg.Xkcd.Parallel)
		return nil
	}

	// Создаем канал для ограничения параллельных запросов
	semaphore := make(chan struct{}, a.cfg.Xkcd.Parallel)

	// Канал для передачи ошибок из горутин
	errCh := make(chan error)

	// Канал для сигнала о завершении
	done := make(chan struct{})

	// Запускаем цикл в горутине для параллельной обработки
	go func() {
		defer close(done)
		for i := 100; i < 100+a.cfg.Xkcd.EndComics; i++ {
			// Получаем токен из канала семафора
			semaphore <- struct{}{}
			go func(id int) {
				defer func() {
					// Освобождаем токен семафора
					<-semaphore
				}()

				comics, err := a.GetComicsByID(id)
				if err != nil {
					errCh <- fmt.Errorf("error get comics from http: %w", err)
					return
				}

				words := a.stem.stemmedWords(comics)

				d := types.Comics{
					Keywords:    words,
					URL:         comics.URL,
					PreKeywords: "",
				}

				// Добавляем данные в хранилище
				a.store.Add(d, id)
			}(i)
		}
	}()

	// Ожидаем завершения всех горутин
	go func() {
		<-done
		close(errCh)
	}()

	// Ждем, пока все горутины завершат выполнение или произойдет ошибка
	for err := range errCh {
		if err != nil {
			return err
		}
	}

	// Сохраняем данные после завершения обработки всех горутин
	if err := a.store.Save(); err != nil {
		return err
	}

	return nil
}
