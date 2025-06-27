-- name: GetFilteredEvents :many
SELECT *
FROM events
WHERE events.user_id = @user_id
  AND (
      @start_date::timestamptz IS NULL OR
      events.start_date >= @start_date
  )
  AND (
      @end_date::timestamptz IS NULL OR
      events.start_date < @end_date
  )
  AND (
      @tag::text IS NULL OR
      EXISTS (
          SELECT 1
          FROM event_tags et
          JOIN tags t ON t.id = et.tag_id
          WHERE et.event_id = events.id AND t.name = @tag
      )
  );
