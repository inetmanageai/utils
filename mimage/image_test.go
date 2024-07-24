package mimage_test

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/inetmanageai/utils/mimage"
	"github.com/stretchr/testify/assert"
)

func createTestImage(types string) []byte {
	img := image.NewRGBA(image.Rect(0, 0, 200, 200))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	buf := new(bytes.Buffer)
	switch types {
	case "jpeg":
		jpeg.Encode(buf, img, nil)
	case "png":
		png.Encode(buf, img)
	}
	return buf.Bytes()
}

func TestPlotImageFromUrl(t *testing.T) {
	type Input struct {
		URL      string
		PlotData []mimage.PlotDataModel
	}
	tests := []struct {
		Name          string
		Input         Input
		ExpectedError bool
	}{
		// TODO: Add test cases.
		{
			Name: "Valid PNG URL and PlotData",
			Input: Input{
				URL: "https://freshfeeds.in/img_urls/img16.png",
				PlotData: []mimage.PlotDataModel{
					{
						Rect:  image.Rect(10, 10, 50, 50),
						Label: "Test Label 1",
					},
					{
						Rect:  image.Rect(0, 0, 5, 5),
						Label: "Test Label 2",
					},
				},
			},
			ExpectedError: false,
		},
		{
			Name: "Valid JPEG URL and PlotData",
			Input: Input{
				URL: "https://picsum.photos/200.jpg",
				PlotData: []mimage.PlotDataModel{
					{
						Rect:  image.Rect(10, 10, 50, 50),
						Label: "Test Label 1",
					},
					{
						Rect:  image.Rect(0, 0, 5, 5),
						Label: "Test Label 2",
					},
				},
			},
			ExpectedError: false,
		},
		{
			Name: "Invalid URL",
			Input: Input{
				URL: "https://picsums.photos/200",
				PlotData: []mimage.PlotDataModel{
					{
						Rect:  image.Rect(10, 10, 50, 50),
						Label: "Test Label",
					},
				},
			},
			ExpectedError: true,
		},
		{
			Name: "Invalid URL type",
			Input: Input{
				URL: "https://picsum.photos/200.webp",
				PlotData: []mimage.PlotDataModel{
					{
						Rect:  image.Rect(10, 10, 50, 50),
						Label: "Test Label 1",
					},
				},
			},
			ExpectedError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// --------------- Act ---------------
			result, err := mimage.PlotImageFromUrl(tt.Input.URL, tt.Input.PlotData)

			// --------------- Assert ---------------
			if tt.ExpectedError {
				assert.Error(t, err)
				assert.Zero(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotZero(t, result)
				_, typeImgResult, _ := image.Decode(bytes.NewBuffer(result))
				res, err := http.Get(tt.Input.URL)
				if err != nil {
					t.Fatal(err)
				}
				defer res.Body.Close()
				f, _ := io.ReadAll(res.Body)
				_, typeImgSrc, _ := image.Decode(bytes.NewBuffer(f))
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, typeImgSrc, typeImgResult)
			}
		})
	}
}

func TestPlotImageFromByte(t *testing.T) {
	type Input struct {
		Byte     []byte
		PlotData []mimage.PlotDataModel
	}
	tests := []struct {
		Name          string
		Input         Input
		ExpectedError bool
	}{
		// TODO: Add test cases.
		{
			Name: "Valid PNG image data with plot data",
			Input: Input{
				Byte: createTestImage("png"),
				PlotData: []mimage.PlotDataModel{
					{
						Rect:  image.Rect(10, 10, 50, 50),
						Label: "Test Label 1",
					},
					{
						Rect:  image.Rect(0, 0, 5, 5),
						Label: "Test Label 2",
					},
				},
			},
			ExpectedError: false,
		},
		{
			Name: "Valid JPEG image data with plot data",
			Input: Input{
				Byte: createTestImage("jpeg"),
				PlotData: []mimage.PlotDataModel{
					{
						Rect:  image.Rect(10, 10, 50, 50),
						Label: "Test Label 1",
					},
					{
						Rect:  image.Rect(0, 0, 5, 5),
						Label: "Test Label 2",
					},
				},
			},
			ExpectedError: false,
		},
		{
			Name: "Invalid image data",
			Input: Input{
				Byte: []byte("invalid data"),
				PlotData: []mimage.PlotDataModel{
					{
						Rect:  image.Rect(10, 10, 50, 50),
						Label: "Test Label",
					},
				},
			},
			ExpectedError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// --------------- Act ---------------
			result, err := mimage.PlotImageFromBytes(tt.Input.Byte, tt.Input.PlotData)

			// --------------- Assert ---------------
			if tt.ExpectedError {
				assert.Error(t, err)
				assert.Zero(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotZero(t, result)
				_, typeImgResult, _ := image.Decode(bytes.NewBuffer(result))
				_, typeImgSrc, _ := image.Decode(bytes.NewBuffer(tt.Input.Byte))
				assert.Equal(t, typeImgSrc, typeImgResult)
			}
		})
	}
}

func TestPlotImageFromDir(t *testing.T) {
	type Input struct {
		FilePath string
		PlotData []mimage.PlotDataModel
	}
	tests := []struct {
		Name          string
		Input         Input
		ExpectedError bool
	}{
		// TODO: Add test cases.
		{
			Name: "Valid PNG image data with plot data",
			Input: Input{
				FilePath: "../testdata/image_test.png",
				PlotData: []mimage.PlotDataModel{
					{
						Rect:  image.Rect(10, 10, 50, 50),
						Label: "Test Label 1",
					},
					{
						Rect:  image.Rect(0, 0, 5, 5),
						Label: "Test Label 2",
					},
				},
			},
			ExpectedError: false,
		},
		{
			Name: "Valid JPEG image data with plot data",
			Input: Input{
				FilePath: "../testdata/image_test.jpg",
				PlotData: []mimage.PlotDataModel{
					{
						Rect:  image.Rect(10, 10, 50, 50),
						Label: "Test Label 1",
					},
					{
						Rect:  image.Rect(0, 0, 5, 5),
						Label: "Test Label 2",
					},
				},
			},
			ExpectedError: false,
		},
		{
			Name: "Invalid image data",
			Input: Input{
				FilePath: "xxx.jpg",
				PlotData: []mimage.PlotDataModel{
					{
						Rect:  image.Rect(10, 10, 50, 50),
						Label: "Test Label",
					},
				},
			},
			ExpectedError: true,
		},
		{
			Name: "Invalid image type",
			Input: Input{
				FilePath: "../testdata/image_test.webp",
				PlotData: []mimage.PlotDataModel{
					{
						Rect:  image.Rect(10, 10, 50, 50),
						Label: "Test Label 1",
					},
				},
			},
			ExpectedError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// --------------- Act ---------------
			result, err := mimage.PlotImageFromDir(tt.Input.FilePath, tt.Input.PlotData)

			// --------------- Assert ---------------
			if tt.ExpectedError {
				assert.Error(t, err)
				assert.Zero(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotZero(t, result)
				_, typeImgResult, _ := image.Decode(bytes.NewBuffer(result))
				f, err := os.Open(tt.Input.FilePath)
				if err != nil {
					t.Fatal(err)
				}
				defer f.Close()
				_, typeImgSrc, _ := image.Decode(f)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, typeImgSrc, typeImgResult)
			}
		})
	}
}
