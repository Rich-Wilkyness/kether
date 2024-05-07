package models

import "time"

// Both teachers and students will be users, and access to classes will be determined by the class teacherID
// Access level is really just for User vs Admin of the site
type User struct {
	ID          int // probably index this
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// UserID is the teacher who created the test
// ClassID is the class that the test is for
type Test struct {
	ID        int // probably index this
	Name      string
	Version   string
	ClassID   int
	UserID    int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TestID is the test that the question belongs to. can be Null if not assigned to a test yet. downside is that a question can only belong to one test, so we would need to duplicate questions for multiple tests
// UserID is the teacher who created the question
// type is multiple choice, true/false, fill in the blank, etc. This allows for different formatting
// for true/false questions, we can have a single correct answer defined in the option table rather than multple true or false answers littered in the table
// use radio buttons to select type
type Question struct {
	ID               int
	Type             string
	Question         string
	DifficultyLevel  int
	TestID           int // probably index this
	UserID           int
	TopicSelectionID int
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// By having the correct field in Option, this allows us to have multiple correct answers defined for a subject easily
type AnswerChoices struct {
	ID         int // probably index this
	Choice     string
	Correct    bool
	QuestionID int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// classes are just used to group curriculum together, it is not inherently specific to grade or subject
// example: parent creates a class for wilderness survival - not grade specific
type Class struct {
	ID        int // probably index this
	TeacherID int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// not sure if this is what i want to do. this could create problems of drop down menues loading slowly and might not be great for filtering. not sure how efficient this method is.
// Admins will generate Grades, Subjects, and Topics based on nationally defined curriculum subjects.
// example: GradeSubject Grade: 1, Subject: Math, Topic: Adding
// goal: to provide a list of topics for teachers to choose from when creating questions to help direct coverage of material
type TopicSelection struct {
	ID        int
	Grade     int
	Subject   string
	Topic     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// this allows us to assign a user to multiple classes because a user can have multiple UserClass entries
type UserInClass struct {
	ID        int
	UserID    int
	ClassID   int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Depending on the scale of your application, consider adding indexes on frequently queried columns to improve query performance.
// for example I if a user needs to take a test, they will click on a link with a Test.ID, this will query the database for Question.TestID = Test.ID
// if you have a lot of questions, you may want to index the TestID column in the Question table to speed up the query

// downside of indexing is that it can slow down write operations, so only index columns that are frequently queried and not updated, inserted, or deleted often
// in the example above, the TestID column in the Question table would be a good candidate for indexing because it is queried often and not updated often
// a user could update the Question.Question field to change how they want the question to be worded, but the TestID would remain the same
