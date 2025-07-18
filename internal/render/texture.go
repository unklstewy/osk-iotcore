package render

import (
	"image"
	"image/draw"
)

// Texture represents a texture that can be rendered
type Texture struct {
	ID     uint32
	Width  int32
	Height int32
	Data   []byte
}

// TextureManager manages texture loading and caching
type TextureManager struct {
	textures map[string]*Texture
}

// NewTextureManager creates a new texture manager
func NewTextureManager() *TextureManager {
	return &TextureManager{
		textures: make(map[string]*Texture),
	}
}

// LoadTexture loads a texture from an image
func (tm *TextureManager) LoadTexture(name string, img image.Image) (*Texture, error) {
	// Convert image to RGBA if needed
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Src)

	texture := &Texture{
		Width:  int32(rgba.Bounds().Dx()),
		Height: int32(rgba.Bounds().Dy()),
		Data:   rgba.Pix,
	}

	tm.textures[name] = texture
	return texture, nil
}

// GetTexture retrieves a cached texture by name
func (tm *TextureManager) GetTexture(name string) (*Texture, bool) {
	tex, exists := tm.textures[name]
	return tex, exists
}

// RenderTexture renders a texture at the specified coordinates
func (tm *TextureManager) RenderTexture(renderer Renderer, texture *Texture, x, y, width, height int) error {
	// Texture rendering logic would go here
	// This is a placeholder that would be implemented based on the specific renderer
	return nil
}

// UnloadTexture removes a texture from memory
func (tm *TextureManager) UnloadTexture(name string) {
	delete(tm.textures, name)
}

// Clear removes all textures from memory
func (tm *TextureManager) Clear() {
	tm.textures = make(map[string]*Texture)
}
