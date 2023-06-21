package queries

const (
	QuerySaveAge = `
	INSERT INTO
    users (chatid, age)
	VALUES
    (
		$1,
        $2
    )
	RETURNING chatid, age;
	`

	QuerySaveGender = `
	UPDATE users 
	SET ismale = $2
	WHERE chatid = $1;
	`

	QuerySaveNationality = `
	UPDATE users 
	SET nationality = $2
	WHERE chatid = $1;
	`

	QuerySelectParameters = `
	SELECT age, nationality
	FROM users
	WHERE chatid = $1;
	`
)