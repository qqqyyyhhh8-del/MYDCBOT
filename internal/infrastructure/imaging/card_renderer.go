package imaging

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	_ "image/jpeg"
	"os"
	"path/filepath"

	"github.com/user/dcminigames/internal/domain/uno/entity"
)

type CardRenderer struct {
	assetsPath string
}

func NewCardRenderer(assetsPath string) *CardRenderer {
	return &CardRenderer{assetsPath: assetsPath}
}

func (r *CardRenderer) RenderHand(cards []*entity.Card) ([]byte, error) {
	if len(cards) == 0 {
		return nil, fmt.Errorf("没有卡牌")
	}

	images := make([]image.Image, 0, len(cards))
	for _, card := range cards {
		img, err := r.loadCardImage(card)
		if err != nil {
			return nil, err
		}
		images = append(images, img)
	}

	bounds := images[0].Bounds()
	cardW, cardH := bounds.Dx(), bounds.Dy()
	overlap := cardW / 3
	totalW := cardW + (len(cards)-1)*(cardW-overlap)

	canvas := image.NewRGBA(image.Rect(0, 0, totalW, cardH))
	for i, img := range images {
		x := i * (cardW - overlap)
		draw.Draw(canvas, image.Rect(x, 0, x+cardW, cardH), img, bounds.Min, draw.Over)
	}

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, canvas, &jpeg.Options{Quality: 85}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (r *CardRenderer) RenderSingleCard(card *entity.Card) ([]byte, error) {
	img, err := r.loadCardImage(card)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 85}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (r *CardRenderer) loadCardImage(card *entity.Card) (image.Image, error) {
	path := filepath.Join(r.assetsPath, card.ImageKey)
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("打开图片失败 %s: %w", card.ImageKey, err)
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("解码图片失败: %w", err)
	}
	return img, nil
}
