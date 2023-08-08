package messenger

import (
	"fmt"
	"log"
	"sync"

	"github.com/Masterminds/squirrel"
	"github.com/gocql/gocql"
	"github.com/timdevlet/todo/internal/helpers"
	"github.com/timdevlet/todo/pkg/postgres"
)

type MessageRepository struct {
	postgres *postgres.PDB

	directChennelMap sync.Map
	session          *gocql.Session
}

func NewMessageRepository(db *postgres.PDB) *MessageRepository {
	cluster := gocql.NewCluster("cassandra", "cassandra", "cassandra")
	cluster.Keyspace = "m"

	session, err := cluster.CreateSession()

	if err != nil {
		log.Fatal(err)
	}

	return &MessageRepository{
		postgres:         db,
		directChennelMap: sync.Map{},
		session:          session,
	}
}

// ----------------------------

func (repo *MessageRepository) InsertMessage(from_uuid, to_uuid string, text ...string) (string, error) {

	channelUuid := helpers.UuidByTwoStrings(from_uuid, to_uuid)

	tx := repo.postgres.DB

	if from_uuid < to_uuid {
		from_uuid, to_uuid = to_uuid, from_uuid
	}

	{
		s1 := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
			Insert("channels_direct").
			Columns("uuid", "from_uuid", "to_uuid", "messages").
			Values(channelUuid, from_uuid, to_uuid, 1).
			Suffix("ON CONFLICT(uuid) DO UPDATE SET updated_at = NOW(), messages = channels_direct.messages + 1")

		sql, values, _ := s1.ToSql()
		_, err := tx.Exec(sql, values...)
		if err != nil {
			return channelUuid, fmt.Errorf("[messenger][repo:InsertMessage] %w -- %s", err, sql)
		}
	}

	{
		s2 := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
			Insert("messages2").
			Columns("channel_uuid", "owner_uuid", "text").
			Suffix("RETURNING \"uuid\"")

		for _, t := range text {
			s2 = s2.Values(channelUuid, from_uuid, t)
		}

		sql, values, _ := s2.ToSql()
		_, err := tx.Exec(sql, values...)
		if err != nil {
			return channelUuid, err
		}
	}

	// g.Go(func() error {

	// 	for _, t := range text {
	// 		if err := repo.session.Query(`INSERT INTO messages (channel_uuid, author_uuid, message, id) VALUES (?, ?, ?, ?)`,
	// 			channelUuid, from_uuid, t, gocql.TimeUUID()).Exec(); err != nil {
	// 			log.Fatal(err)
	// 		}
	// 	}

	// 	return nil

	// })

	return channelUuid, nil
}

func (repo *MessageRepository) FetchMessagesByUser(from_uuid string) ([]Message, error) {

	statement := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("uuid", "from_uuid", "to_uuid").
		From("channels_direct").
		Where(squirrel.Or{squirrel.Eq{"from_uuid": from_uuid}, squirrel.Eq{"to_uuid": from_uuid}}).
		Limit(100)

	sql, values, _ := statement.ToSql()

	rows, err := repo.postgres.DB.Query(sql, values...)
	if err != nil {
		return []Message{}, err
	}
	defer rows.Close()

	result := make([]ChannelDirect, 0)
	for rows.Next() {
		var i ChannelDirect
		err := rows.Scan(&i.UUID, &i.FromUuid, &i.ToUuid)
		if err != nil {
			return []Message{}, err
		}
		result = append(result, i)
	}

	if len(result) == 0 {
		return []Message{}, nil
	}

	if err = rows.Err(); err != nil {
		return []Message{}, err
	}

	channelDirect := result[0]

	return repo.FetchMessages(channelDirect.UUID)
}

func (repo *MessageRepository) FetchMessages(channel_uuid string) ([]Message, error) {
	statement := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("uuid, owner_uuid, channel_uuid, text").
		From("messages2").
		Where(squirrel.Eq{"channel_uuid": channel_uuid}).
		Limit(uint64(100000))

	sql, values, _ := statement.ToSql()

	rows, err := repo.postgres.DB.Query(sql, values...)
	if err != nil {
		return []Message{}, err
	}
	defer rows.Close()

	result := make([]Message, 0)
	for rows.Next() {
		var i Message
		err := rows.Scan(&i.UUID, &i.OwnerUuid, &i.ChannelUuid, &i.Text)
		if err != nil {
			return []Message{}, err
		}
		result = append(result, i)
	}

	if err = rows.Err(); err != nil {
		return []Message{}, err
	}

	return result, nil
}

func (repo *MessageRepository) FetchChannels(user_uuid string) ([]ChannelDirect, error) {

	statement := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("uuid", "from_uuid", "to_uuid", "created_at").
		From("channels_direct").
		Where(squirrel.Or{squirrel.Eq{"from_uuid": user_uuid}, squirrel.Eq{"to_uuid": user_uuid}}).
		Limit(1)

	sql, values, _ := statement.ToSql()

	rows, err := repo.postgres.DB.Query(sql, values...)
	if err != nil {
		return []ChannelDirect{}, err
	}
	defer rows.Close()

	result := make([]ChannelDirect, 0)
	for rows.Next() {
		var i ChannelDirect
		err := rows.Scan(&i.UUID, &i.FromUuid, &i.ToUuid, &i.CreatedAt)
		if err != nil {
			return []ChannelDirect{}, err
		}
		result = append(result, i)
	}

	if err = rows.Err(); err != nil {
		return []ChannelDirect{}, err
	}

	return result, nil
}

func (repo *MessageRepository) CountChannelMessages(channel_uuid string) (int, error) {

	statement := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("count(*) as total").
		From("messages2").
		Where(squirrel.Eq{"channel_uuid": channel_uuid}).
		Limit(1)

	sql, values, _ := statement.ToSql()

	rows, err := repo.postgres.DB.Query(sql, values...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var total int
	rows.Next()

	err = rows.Scan(&total)

	if err != nil {
		return 0, err
	}

	return total, nil
}
