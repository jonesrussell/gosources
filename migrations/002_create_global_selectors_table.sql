-- Create global_selectors table for default selector configuration
CREATE TABLE IF NOT EXISTS global_selectors (
    id SERIAL PRIMARY KEY,
    selectors JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create trigger to automatically update updated_at
CREATE TRIGGER update_global_selectors_updated_at
    BEFORE UPDATE ON global_selectors
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Insert default global selectors with comprehensive patterns
INSERT INTO global_selectors (selectors) VALUES (
    '{
        "article": {
            "container": "article, main, .article-content, [itemprop=\"articleBody\"]",
            "title": "h1, h2.title, meta[property=\"og:title\"]",
            "body": "[itemprop=\"articleBody\"], .article-body, .article-content, article",
            "intro": "meta[property=\"og:description\"], meta[name=\"description\"], .article-intro, .lead",
            "link": "a[href*=\"/news/\"], a[href*=\"/article/\"], .article-link a",
            "image": "meta[property=\"og:image\"], img[itemprop=\"image\"], article img, .article-image img",
            "byline": "[itemprop=\"author\"], .author, .byline, meta[name=\"author\"]",
            "published_time": "time[datetime], meta[property=\"article:published_time\"], [itemprop=\"datePublished\"]",
            "time_ago": "time, .published-time, .article-time",
            "section": "meta[property=\"article:section\"], .category, .section",
            "category": "meta[property=\"article:section\"], .category, .section-name",
            "article_id": "[data-article-id], article[id]",
            "json_ld": "script[type=\"application/ld+json\"]",
            "keywords": "meta[name=\"keywords\"]",
            "description": "meta[name=\"description\"]",
            "og_title": "meta[property=\"og:title\"]",
            "og_description": "meta[property=\"og:description\"]",
            "og_image": "meta[property=\"og:image\"]",
            "og_url": "meta[property=\"og:url\"]",
            "og_type": "meta[property=\"og:type\"]",
            "og_site_name": "meta[property=\"og:site_name\"]",
            "canonical": "link[rel=\"canonical\"]",
            "author": "meta[property=\"article:author\"], meta[name=\"author\"], [itemprop=\"author\"]",
            "exclude": [
                ".ad, .advertisement, [class*=\"ad-\"], [id^=\"ad-\"], [data-ad], .sponsored",
                "nav, .header, .footer, .navigation, .menu, .site-header, .site-footer",
                "script, style, noscript, [aria-hidden=\"true\"], .visually-hidden, .hidden, svg",
                ".widget, .sidebar, .related-articles, .related-posts",
                ".video-player, .weather, .newsletter",
                ".social-share, .share-buttons, button, form, .subscribe",
                ".comments, .comment-section, .user-comments",
                ".pagination, .breadcrumb, .tags, .author-bio, .more-stories"
            ]
        },
        "list": {
            "container": ".article-list, .news-feed, main, .content",
            "article_cards": "article, .article-item, .news-item",
            "article_list": ".articles, .news-list, ol, ul.articles",
            "exclude_from_list": [
                ".ad, .advertisement, .sponsored, [data-ad]"
            ]
        },
        "page": {
            "container": "main, article, .content, .main-content, #main-content, .page-content",
            "title": "h1, meta[property=\"og:title\"], .page-title",
            "content": "main, article, .content, .page-content, .main-content",
            "description": "meta[name=\"description\"]",
            "keywords": "meta[name=\"keywords\"]",
            "og_title": "meta[property=\"og:title\"]",
            "og_description": "meta[property=\"og:description\"]",
            "og_image": "meta[property=\"og:image\"]",
            "og_url": "meta[property=\"og:url\"]",
            "canonical": "link[rel=\"canonical\"]",
            "exclude": [
                ".ad, .advertisement, [class*=\"ad-\"], [id^=\"ad-\"], [data-ad], .sponsored",
                "nav, .header, .footer, .navigation, .menu, .site-header, .site-footer, .breadcrumb",
                "script, style, noscript, [aria-hidden=\"true\"], .visually-hidden, .hidden, svg",
                ".widget, .sidebar, .video-player, .weather, .newsletter",
                ".social-share, .share-buttons, button, form, .subscribe",
                ".pagination, .comments, .author-bio, .more-stories, .related-posts",
                "article.article-item, .article-card, .news-item"
            ]
        }
    }'::jsonb
);
