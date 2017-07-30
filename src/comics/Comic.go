package comics


type Comic struct {
	title string
	writers []string
	published string
	imgPath string
	price float32
	src string
}

func NewComic(title string, writers []string, published string, imgPath string, price float32, src string) *Comic {
	return &Comic {
		title:   title,
		writers: writers,
		published: published,
		imgPath: imgPath,
		price:   price,
		src:     src,
	}
}

func (c *Comic) Save() bool {
	// @todo: Implement
	return false
}

func (c *Comic) GetTitle() string {
	return c.title
}

func (c *Comic) GetWriters() []string {
	return c.writers
}

func (c *Comic) GetImgPath() string {
	return c.imgPath
}

func (c *Comic) GetPublished() string {
	return c.published
}

func (c *Comic) GetPrice() float32 {
	return c.price
}

func (c *Comic) GetSource() string {
	return c.src
}