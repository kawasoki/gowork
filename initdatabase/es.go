package initdatabase

import (
	"context"
	"fmt"
	"github.com/kawasoki/gowork/configs"
	"github.com/kawasoki/gowork/logger"
	"strings"
	"sync"
	"time"

	"github.com/olivere/elastic/v7"
)

var (
	esClient    *elastic.Client
	mutex       sync.Mutex
	esProcessor *EsBulkProcessor
)

func NewEsClient(config *configs.Config) {
	address := strings.Split(config.EsAddress, ",")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	opts := []elastic.ClientOptionFunc{
		elastic.SetURL(address...),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(true),
	}
	if config.EsUserName != "" || config.EsPassword != "" {
		opts = append(opts, elastic.SetBasicAuth(config.EsUserName, config.EsPassword))
	}
	client, err := elastic.DialContext(ctx, opts...)
	if err != nil {
		panic(fmt.Errorf("无法连接es: %v", err))
	}
	esClient = client
}

func GetEs() {
	NewEsClient(&configs.Config{})
}

var once sync.Once

func GetEsProcessor() *EsBulkProcessor {
	once.Do(
		func() {
			client := GetEs()
			p, err := NewEsBulkProcessor(client, 2000, 4, 60)
			if err != nil {
				panic("NewEsBulkProcessor panic" + err.Error())
			}
			esProcessor = p
		})
	return esProcessor
}

type EsBulkProcessor struct {
	Client        *elastic.Client
	BulkActions   int
	Worker        int
	FlushTime     int
	BulkProcessor *elastic.BulkProcessor
}

func NewEsBulkProcessor(client *elastic.Client, bulkActions int, worker int, flushTime int) (*EsBulkProcessor, error) {
	p := &EsBulkProcessor{client, bulkActions, worker, flushTime, nil}
	processor, err := p.Client.BulkProcessor().BulkActions(p.BulkActions).FlushInterval(time.Second * time.Duration(p.FlushTime)).After(p.getFailed).Workers(p.Worker).Do(context.TODO())
	if err != nil {
		panic("初始化es批量进程失败")
	}
	p.BulkProcessor = processor
	err = p.BulkProcessor.Start(context.Background())
	return p, err
}
func (p *EsBulkProcessor) getFailed(executionId int64, requests []elastic.BulkableRequest, response *elastic.BulkResponse, err error) {
	if response == nil {
		return
	}
	fi := response.Failed()
	if len(fi) != 0 {
		for _, f := range fi {
			logger.Infof("DebugFailedEs: index:%s type:%s id:%s version:%d  status:%d result:%s ForceRefresh:%v errorDetail:%v getResult:%v", f.Index, f.Type, f.Id, f.Version, f.Status, f.Result, f.ForcedRefresh, f.Error, f.GetResult)
		}
	}
}

func (p *EsBulkProcessor) BatchInsert(index string, data interface{}) {
	request := elastic.NewBulkIndexRequest().Index(index).Doc(data)
	p.BulkProcessor.Add(request)

}

func (p *EsBulkProcessor) BatchInsertWithId(id, index string, data interface{}) {
	request := elastic.NewBulkIndexRequest().Id(id).Index(index).Doc(data)
	p.BulkProcessor.Add(request)
}
