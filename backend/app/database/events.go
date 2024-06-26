package database

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type DBEvent struct {
	Id              pgtype.UUID      `db:"id"`
	StatusId        int              `db:"status_id"`
	Title           string           `db:"title"`
	Description     string           `db:"description"`
	Rules           string           `db:"rules"`
	Timeline        string           `db:"timeline"`
	OrganizerUserId pgtype.UUID      `db:"organizer_user_id" json:"-"`
	MaxTeams        int              `db:"max_teams"`
	StartsAt        pgtype.Timestamp `db:"starts_at"`
	EndsAt          pgtype.Timestamp `db:"ends_at"`
	CreatedOn       pgtype.Timestamp `db:"created_on" json:"-"`
}

type DBEventStatus struct {
	Id          int    `db:"id"`
	Code        string `db:"code"`
	Title       string `db:"title"`
	Description string `db:"description"`
	Sort        int    `db:"sort"`
}

type DBEventStatusCode struct {
	Code string `db:"code"`
}

func GetEventStatusCode(eventId pgtype.UUID) (string, error) {
	result, err := GetRow[DBEventStatusCode](
		`SELECT code
         FROM events 
         INNER JOIN statuses ON (events.status_id = statuses.id)
         WHERE events.id = $1`,
		eventId)
	return result.Code, err
}

func CreateEvent(organizerUserId pgtype.UUID) (DBEvent, error) {
	var event DBEvent
	event, err := GetRow[DBEvent](
		`INSERT INTO events
            (status_id, title, description, rules, organizer_user_id)
            VALUES
            ((SELECT id FROM statuses WHERE code = 'PLANNING'),
            '',
            '',
            '',
            $1)
         RETURNING *
        `,
		organizerUserId)
	return event, err
}

func GetEvent(eventId pgtype.UUID) (DBEvent, error) {
	event, err := GetRow[DBEvent](
		`SELECT * FROM events WHERE id = $1`,
		eventId)
	return event, err
}

func GetEvents() ([]DBEvent, error) {
	result, err := GetRows[DBEvent](`SELECT * FROM events`)
	return result, err
}

// GetActiveEvent will return what is assumed to be a single active event.
func GetActiveEvent() (DBEvent, error) {
	result, err := GetRow[DBEvent](
		`SELECT * FROM events
         WHERE status_id IN (
           SELECT id FROM statuses WHERE code NOT IN ('PLANNING', 'ENDED') 
         )
         LIMIT 1`)
	return result, err
}

func UpdateEvent(event DBEvent) (DBEvent, error) {
	event, err := GetRow[DBEvent](
		`UPDATE events
         SET status_id=$2,
             title=$3,
             timeline=$4,
             description=$5,
             rules=$6,
             max_teams=$7,
             starts_at=$8,
             ends_at=$9
         WHERE id=$1
         RETURNING *`,
		event.Id, event.StatusId, event.Title, event.Timeline, event.Description, event.Rules, event.MaxTeams, event.StartsAt, event.EndsAt)
	return event, err
}

func GetStatuses() ([]DBEventStatus, error) {
	statuses, err := GetRows[DBEventStatus](`SELECT * FROM statuses ORDER BY sort`)
	return statuses, err
}
