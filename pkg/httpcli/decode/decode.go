package decode

import (
	"mime"
	"sync"

	"github.com/wangweihong/gotoolbox/pkg/errors"
)

const (
	ApplicationXml  = "application/xml"
	ApplicationJson = "application/json"
	ContentType     = "Content-Type"
)

// UnmarshalFunc implements manifest unmarshalling a given MediaType.
type UnmarshalFunc func([]byte, interface{}) error

// NewMarshalMapping create a new MarshalMapping.
func NewMarshalMapping() *MarshalMapping {
	return &MarshalMapping{
		mappings: make(map[string]UnmarshalFunc),
	}
}

type MarshalMapping struct {
	lock     sync.RWMutex
	mappings map[string]UnmarshalFunc
}

// ManifestMediaTypes returns the supported media types for manifests.
func (mm *MarshalMapping) ManifestMediaTypes() (mediaTypes []string) {
	mm.lock.RLock()
	defer mm.lock.RUnlock()

	for t := range mm.mappings {
		if t != "" {
			mediaTypes = append(mediaTypes, t)
		}
	}
	return
}

// ManifestMediaTypes returns the supported media types for manifests.
func (mm *MarshalMapping) Register(mediaType string, u UnmarshalFunc) error {
	mm.lock.Lock()
	defer mm.lock.Unlock()

	if _, ok := mm.mappings[mediaType]; ok {
		return errors.Errorf("manifest media type registration would overwrite existing: %s", mediaType)
	}
	mm.mappings[mediaType] = u
	return nil
}

// UnmarshalManifest looks up manifest unmarshal functions based on MediaType.
func (mm *MarshalMapping) UnmarshalManifest(ctHeader string, p []byte, arg interface{}) error {
	mm.lock.Lock()
	defer mm.lock.Unlock()

	var mediaType string
	if ctHeader != "" {
		var err error
		// 相对于直接读取Content-Type
		// mime.ParseMediaType会解析Content-Type并拆分成主类型和可选参数
		// "application/json;
		// charset=utf-8",mime.ParseMediaType将返回contentType为"application/json"，params为map[string]string{"charset":
		// "utf-8"}
		mediaType, _, err = mime.ParseMediaType(ctHeader)
		if err != nil {
			return err
		}
	}

	unmarshalFunc, ok := mm.mappings[mediaType]
	if !ok {
		return errors.Errorf("unsupported Content-Type %v", mediaType)
	}

	return unmarshalFunc(p, arg)
}
