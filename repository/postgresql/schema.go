package postgresql

var migrate = []string{
	`
    CREATE TABLE IF NOT EXISTS "user" (
      id            BIGSERIAL    NOT NULL PRIMARY KEY,
      username      TEXT         NOT NULL,
      firstname     TEXT         NOT NULL,
      lastname      TEXT         NOT NULL,
      passwordhash  TEXT         NOT NULL,
      createdat     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
      updatedat     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
      admin         BOOLEAN      NOT NULL DEFAULT FALSE
    );
    CREATE UNIQUE INDEX IF NOT EXISTS unique_username ON "user" (LOWER(username));
  `,
	`
    CREATE TABLE IF NOT EXISTS problem (
      id               BIGSERIAL    NOT NULL PRIMARY KEY,
      createdat        TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
      updatedat        TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
      authorid         BIGINT       NOT NULL REFERENCES "user"(id),
      name             TEXT         NOT NULL,
      externalurl      TEXT         NOT NULL
    );
  `,
	`
    CREATE TABLE IF NOT EXISTS problem_revision (
      id               BIGSERIAL    NOT NULL PRIMARY KEY,
      createdat        TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
      updatedat        TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
      authorid         BIGINT       NOT NULL REFERENCES "user"(id),
      problemid        BIGINT       NOT NULL REFERENCES "problem"(id),
      revision         INTEGER      NOT NULL,
      title            TEXT         NOT NULL,
      memorylimit      INTEGER      NOT NULL,
      timelimit        INTEGER      NOT NULL,
      statementhtml    TEXT         NOT NULL,
      statementpdf     BYTEA        NOT NULL,
      packagearchive   BYTEA        NOT NULL,
      outdated         BOOL         NOT NULL
    );
  `,
	`
    CREATE TABLE IF NOT EXISTS contest (
      id               BIGSERIAL    NOT NULL PRIMARY KEY,
      createdat        TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
      updatedat        TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
      authorid         BIGINT       NOT NULL REFERENCES "user"(id),
      contestname      TEXT         NOT NULL,
      starttime        TIMESTAMPTZ  NOT NULL,
      endtime          TIMESTAMPTZ  NOT NULL
    );
      `,
	`
    CREATE TABLE IF NOT EXISTS contest_problem (
      contestid          BIGINT        NOT NULL REFERENCES contest(id),
      problemid          BIGINT        NOT NULL REFERENCES problem(id),
      CONSTRAINT contest_problem_pkey  PRIMARY KEY (contestid, problemid)
    );
  `,

	`
  CREATE TABLE IF NOT EXISTS contest_user (
    contestid          BIGINT        NOT NULL REFERENCES contest(id),
    userid             BIGINT        NOT NULL REFERENCES "user"(id),
    CONSTRAINT contest_user_pkey     PRIMARY KEY (contestid, userid)
  );
`,
	`
    CREATE TABLE IF NOT EXISTS submission (
      id               BIGSERIAL     NOT NULL PRIMARY KEY,
      createdat        TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
      updatedat        TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
      problemid        BIGINT        NOT NULL REFERENCES problem(id),
      problemrevid     BIGINT        NOT NULL REFERENCES problem_revision(id),
      contestid        BIGINT        NULL REFERENCES contest(id),
      authorid         BIGINT        NOT NULL REFERENCES "user"(id),
      sourcecode       TEXT          NOT NULL,
      language         TEXT          NOT NULL,
      status           TEXT          NOT NULL,
      result           TEXT          NOT NULL DEFAULT '',
      details          TEXT          NOT NULL DEFAULT ''
    );
  `,
}

var drop = []string{
	`drop table if exists "user" cascade`,
	`drop table if exists problem cascade`,
	`drop table if exists problem_revision cascade`,
	`drop table if exists contest cascade`,
	`drop table if exists contest_problem cascade`,
	`drop table if exists contest_user cascade`,
	`drop table if exists submission cascade`,
}

var seed = []string{
	`
  INSERT INTO public."user" (
    username, firstname, lastname, passwordhash, createdat, admin
  )
  VALUES
    ('root', 'admin', 'admin', '$2a$10$OwjCOmEq7jd5Rc5sg3bOFOwAQmx8/xXx/Mbyt3.2jem.rxs9Imo16',  '2020-05-05 17:17:11.111821+03',  true),
    ('test',  'Test', 'User', '$2a$10$OwjCOmEq7jd5Rc5sg3bOFOwAQmx8/xXx/Mbyt3.2jem.rxs9Imo16',  '2020-05-05 17:17:11.111821+03',  false),
    ('user1',  'John', 'Smith', '$2a$10$OwjCOmEq7jd5Rc5sg3bOFOwAQmx8/xXx/Mbyt3.2jem.rxs9Imo16',  '2020-05-05 17:17:11.111821+03',  false)
  ON CONFLICT DO NOTHING;`,
}
