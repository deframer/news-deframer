SELECT i.id,
       f.url AS feed_url,
       f.root_domain,
       i.media_content
FROM public.items i
JOIN public.feeds f ON f.id = i.feed_id
LIMIT 1000;
