package resources

import (
	"context"
	"strconv"

	"github.com/aviznetworks/ones-gfx-sdk/ones-gfx-sdk-go/ones_gfx"
)

// FilesResource covers artifact upload/list/delete.
type FilesResource struct {
	transport *ones_gfx.Transport
}

// NewFilesResource constructs a FilesResource.
func NewFilesResource(transport *ones_gfx.Transport) *FilesResource {
	return &FilesResource{transport: transport}
}

// UploadFileResult is the response of Upload (POST /uploadfile).
type UploadFileResult struct {
	Success  bool   `json:"success,omitempty"`
	Message  string `json:"message,omitempty"`
	FilePath string `json:"filepath,omitempty"`
	Checksum string `json:"checksum,omitempty"`
	FileSize int64  `json:"filesize,omitempty"`
	FileName string `json:"filename,omitempty"`
	Error    string `json:"error,omitempty"`
}

// GetFilesResult is the response of List (GET /getfiles/{filetype}).
// The Files field is a grouped map (filetype "ALL") or a list (specific
// filetype) depending on the query — surfaced raw for that reason.
type GetFilesResult struct {
	Success    bool        `json:"success,omitempty"`
	Files      interface{} `json:"files,omitempty"`
	TotalCount *int        `json:"totalCount,omitempty"`
	Count      *int        `json:"count,omitempty"`
	Error      string      `json:"error,omitempty"`
}

// DeleteFileResult is the response of Delete (DELETE /deletefile/{id}).
type DeleteFileResult struct {
	Success   bool   `json:"success,omitempty"`
	Message   string `json:"message,omitempty"`
	DeletedID *int   `json:"deletedId,omitempty"`
	Error     string `json:"error,omitempty"`
}

// Upload stores an artifact (NOS image, ONES_T, ONES_CONF, CONFIG...).
// Maps to POST /uploadfile (multipart).
func (r *FilesResource) Upload(ctx context.Context, filePath, filetype string, version, vendor, tag *string) (*UploadFileResult, error) {
	fields := map[string]string{"filetype": filetype}
	if version != nil {
		fields["version"] = *version
	}
	if vendor != nil {
		fields["vendor"] = *vendor
	}
	if tag != nil {
		fields["tag"] = *tag
	}
	return ones_gfx.CallMultipart[UploadFileResult](r.transport, "POST", "uploadfile", fields, map[string]string{"file": filePath}, ones_gfx.OperationModeSynchronous, nil)
}

// List returns uploaded files of one type (or grouped by type for "ALL").
// Maps to GET /getfiles/{filetype}.
func (r *FilesResource) List(ctx context.Context, filetype string) (*GetFilesResult, error) {
	return ones_gfx.Call[GetFilesResult](r.transport, "GET", "getfiles/"+filetype, nil, ones_gfx.OperationModeSynchronous, nil)
}

// Delete removes an uploaded file by ID. Maps to DELETE /deletefile/{id}.
func (r *FilesResource) Delete(ctx context.Context, id int) (*DeleteFileResult, error) {
	return ones_gfx.Call[DeleteFileResult](r.transport, "DELETE", "deletefile/"+strconv.Itoa(id), nil, ones_gfx.OperationModeSynchronous, nil)
}
