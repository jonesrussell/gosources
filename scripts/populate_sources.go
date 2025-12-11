package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	apiBaseURL      = "http://localhost:8050/api/v1/sources"
	sudburySourceID = "96abf22e-2577-4b4b-aea3-c524e5c8dd3a"
)

type Source struct {
	ID           string              `json:"id,omitempty"`
	Name         string              `json:"name"`
	URL          string              `json:"url"`
	ArticleIndex string              `json:"article_index"`
	PageIndex    string              `json:"page_index"`
	RateLimit    string              `json:"rate_limit"`
	MaxDepth     int                 `json:"max_depth"`
	Time         []string            `json:"time,omitempty"`
	Selectors    Selectors           `json:"selectors"`
	CityName     *string             `json:"city_name,omitempty"`
	GroupID      *string             `json:"group_id,omitempty"`
	Enabled      bool                `json:"enabled"`
}

type Selectors struct {
	Article ArticleSelectors `json:"article"`
	List    ListSelectors    `json:"list"`
	Page    PageSelectors    `json:"page"`
}

type ArticleSelectors struct {
	Container      string   `json:"container,omitempty"`
	Title          string   `json:"title,omitempty"`
	Body           string   `json:"body,omitempty"`
	Intro          string   `json:"intro,omitempty"`
	Link           string   `json:"link,omitempty"`
	Image          string   `json:"image,omitempty"`
	Byline         string   `json:"byline,omitempty"`
	PublishedTime  string   `json:"published_time,omitempty"`
	TimeAgo        string   `json:"time_ago,omitempty"`
	Section        string   `json:"section,omitempty"`
	Category       string   `json:"category,omitempty"`
	ArticleID      string   `json:"article_id,omitempty"`
	JSONLD         string   `json:"json_ld,omitempty"`
	Keywords       string   `json:"keywords,omitempty"`
	Description    string   `json:"description,omitempty"`
	OGTitle        string   `json:"og_title,omitempty"`
	OGDescription  string   `json:"og_description,omitempty"`
	OGImage        string   `json:"og_image,omitempty"`
	OGURL          string   `json:"og_url,omitempty"`
	OGType         string   `json:"og_type,omitempty"`
	OGSiteName     string   `json:"og_site_name,omitempty"`
	Canonical      string   `json:"canonical,omitempty"`
	Author         string   `json:"author,omitempty"`
	Exclude        []string `json:"exclude,omitempty"`
}

type ListSelectors struct {
	Container        string   `json:"container,omitempty"`
	ArticleCards     string   `json:"article_cards,omitempty"`
	ArticleList      string   `json:"article_list,omitempty"`
	ExcludeFromList  []string `json:"exclude_from_list,omitempty"`
}

type PageSelectors struct {
	Container      string   `json:"container,omitempty"`
	Title          string   `json:"title,omitempty"`
	Content        string   `json:"content,omitempty"`
	Description    string   `json:"description,omitempty"`
	Keywords       string   `json:"keywords,omitempty"`
	OGTitle        string   `json:"og_title,omitempty"`
	OGDescription  string   `json:"og_description,omitempty"`
	OGImage        string   `json:"og_image,omitempty"`
	OGURL          string   `json:"og_url,omitempty"`
	Canonical      string   `json:"canonical,omitempty"`
	Exclude        []string `json:"exclude,omitempty"`
}

func updateSource(source Source, id string) error {
	data, err := json.Marshal(source)
	if err != nil {
		return fmt.Errorf("marshal source: %w", err)
	}

	url := fmt.Sprintf("%s/%s", apiBaseURL, id)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("update failed with status %d: %s", resp.StatusCode, string(body))
	}

	fmt.Printf("✓ Updated source: %s\n", source.Name)
	return nil
}

func createSource(source Source) error {
	data, err := json.Marshal(source)
	if err != nil {
		return fmt.Errorf("marshal source: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, apiBaseURL, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("create failed with status %d: %s", resp.StatusCode, string(body))
	}

	fmt.Printf("✓ Created source: %s\n", source.Name)
	return nil
}

func main() {
	// Update Sudbury.com source
	sudbury := Source{
		Name:         "sudbury.com",
		URL:          "https://www.sudbury.com/",
		ArticleIndex: "sudbury_com_articles",
		PageIndex:    "sudbury_com_pages",
		RateLimit:    "1s",
		MaxDepth:     2,
		Enabled:      true,
		Selectors: Selectors{
			Article: ArticleSelectors{
				Title:         "meta[property='og:title'], h2.title",
				Body:          "[itemprop='articleBody']",
				PublishedTime: "time[datetime]",
				Image:         "meta[property='og:image']",
				Link:          "a.section-item",
				Category:      "meta[property='article:section'], .section",
				Exclude: []string{
					"nav",
					"script",
					"style",
					"noscript",
					"[aria-hidden='true']",
					"button",
					"form",
					".widget",
				},
			},
			List: ListSelectors{},
			Page: PageSelectors{},
		},
	}

	cityName := "Sudbury"
	groupID := "550e8400-e29b-41d4-a716-446655440000"
	sudbury.CityName = &cityName
	sudbury.GroupID = &groupID

	if err := updateSource(sudbury, sudburySourceID); err != nil {
		fmt.Fprintf(os.Stderr, "Error updating Sudbury source: %v\n", err)
		os.Exit(1)
	}

	// Create Mid-North Monitor source
	midNorth := Source{
		Name:         "Mid-North Monitor",
		URL:          "https://www.midnorthmonitor.com/category/news/local-news/",
		ArticleIndex: "midnorthmonitor_articles",
		PageIndex:    "midnorthmonitor_pages",
		RateLimit:    "1s",
		MaxDepth:     2,
		Time:         []string{"11:45", "23:45"},
		Enabled:      true,
		Selectors: Selectors{
			Article: ArticleSelectors{
				Container:     "article.article-card[data-article-id], main#main-content",
				Title:         "h1, h2.article-card__headline span.article-card__headline-clamp, h3.article-card__headline span.article-card__headline-clamp",
				Body:          ".article-card__excerpt, .article-body, .article-content",
				Intro:         ".article-card__excerpt[data-tb-description], meta[property='og:description']",
				Link:          "a.article-card__link[data-tb-link][href^='/news/'], a[data-tb-link][href*='/news/']",
				Image:         "picture.article-card__image img[data-tb-thumbnail], picture.article-card__image img[data-src], meta[property='og:image']",
				Byline:        ".article-card__meta-bottom .article-author, meta[property='article:author']",
				PublishedTime: "meta[property='article:published_time'], .article-card__time-clamp[data-tb-date]",
				TimeAgo:       ".article-card__time-clamp[data-tb-date], time",
				Section:       ".article-card__category-link[data-tb-category-link] span[data-tb-category]",
				Category:      "meta[property='article:section'], .article-card__category",
				ArticleID:     "[data-article-id]",
				JSONLD:        "script[type='application/ld+json']",
				Keywords:      "meta[name='keywords']",
				Description:   "meta[name='description']",
				OGTitle:       "meta[property='og:title']",
				OGDescription: "meta[property='og:description']",
				OGImage:       "meta[property='og:image']",
				OGURL:         "meta[property='og:url']",
				OGType:        "meta[property='og:type']",
				OGSiteName:    "meta[property='og:site_name']",
				Canonical:     "link[rel='canonical']",
				Author:        "meta[property='article:author']",
				Exclude: []string{
					".ad, .ad__native, [class*='ad__'], [id^='ad-'], [id*='ad__'], [data-aqa='advertisement']",
					"[data-ad], [data-ad-loc], [data-ad-mobile], [data-native]",
					".close-sticky-ad",
					".header, .footer, nav, [class*='-nav'], .header__after",
					"script, style, noscript, [aria-hidden='true'], .visually-hidden, svg",
					"[data-evt], [data-evt-val], [data-evt-typ]",
					"[data-tb-region], [data-tb-region-item]",
					"[data-evt-skip-click]",
					".video-playlist, .video-playlist__queue",
					".weather-widget",
					".widget--local-ads, .local-spotlight",
					".newsletter-widget",
					".list-widget",
					".social-follow, .share-buttons",
					"button, form, .button, .subscribe-btn",
					"img[src*='placeholder']",
					"img[src*='fallback']",
					"img[src*='data:image/svg']",
					"img[src*='quality=5']",
					".ad__placeholder, .placeholder-inner",
					".consent__banner, .fixed-bottom",
					".pagination, .more-stories",
					".related-posts, .hero-feed__widget",
					".sidebar, .comments-section",
					".view-counter, .author-bio",
					".feed-section__content--category li.counter-reset",
					".find-the-best-places",
					"[data-widget-iframe-component]",
					"[data-carousel-item] .ad__native",
				},
			},
			List: ListSelectors{
				Container:    ".hero-feed, .feed-section, [data-tb-region='Main'], [data-tb-region='local-news']",
				ArticleCards: "article.article-card[data-article-id]:not(.ad__native)",
				ArticleList:  "ol.feed-section__content li[data-carousel-item]",
				ExcludeFromList: []string{
					"li.counter-reset",
					".ad__native",
					"article[data-native]",
				},
			},
			Page: PageSelectors{
				Container:     "main, article, .content, .main-content, #main-content",
				Title:         "h1, .page-title, .category-title",
				Content:       "main, article, .content, .page-content",
				Description:   "meta[name='description']",
				Keywords:      "meta[name='keywords']",
				OGTitle:       "meta[property='og:title']",
				OGDescription: "meta[property='og:description']",
				OGImage:       "meta[property='og:image']",
				OGURL:         "meta[property='og:url']",
				Canonical:     "link[rel='canonical']",
				Exclude: []string{
					".ad, .ad__native, [class*='ad__'], [id^='ad-'], [id*='ad__'], [data-aqa='advertisement']",
					"[data-ad], [data-ad-loc], [data-ad-mobile], [data-native]",
					".close-sticky-ad",
					".ad__placeholder, .placeholder-inner",
					".header, .footer, nav, [class*='-nav'], .header__after, .site-header, .site-footer",
					"script, style, noscript, [aria-hidden='true'], .visually-hidden, svg",
					"[data-evt], [data-evt-val], [data-evt-typ], [data-evt-skip-click]",
					"[data-tb-region], [data-tb-region-item]",
					".video-playlist, .video-playlist__queue",
					".weather-widget",
					".widget--local-ads, .local-spotlight",
					".newsletter-widget",
					".list-widget",
					"[data-widget-iframe-component]",
					".social-follow, .share-buttons",
					"button, form, .button, .subscribe-btn",
					".consent__banner, .fixed-bottom",
					".pagination, .more-stories",
					".related-posts, .hero-feed__widget",
					".sidebar, .comments-section",
					".view-counter, .author-bio",
					".find-the-best-places",
					"[data-carousel-item] .ad__native",
					"article.article-card, .feed-section, .hero-feed",
				},
			},
		},
	}

	if err := createSource(midNorth); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating Mid-North Monitor source: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n✓ All sources populated successfully!")
}
