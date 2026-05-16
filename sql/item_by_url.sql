SELECT
    i.id,
    f.url AS feed,
    i.url,
    i.media_content
FROM items i
JOIN feeds f ON f.id = i.feed_id
WHERE i.url = 'https://www.spiegel.de/panorama/justiz/bristol-19-jaehrige-nazi-sympathisantin-fuer-mordversuch-mit-axt-verurteilt-a-285b3edf-c73e-4ed7-97b4-cd1db6a045d8';
