package service

import (
	"testing"
	"time"

	"github.com/yaattc/automatic-time-table-creation/backend/app/store/uni"

	"golang.org/x/crypto/bcrypt"

	"github.com/yaattc/automatic-time-table-creation/backend/app/store/user"

	"github.com/Semior001/timetype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yaattc/automatic-time-table-creation/backend/app/store"
	"github.com/yaattc/automatic-time-table-creation/backend/app/store/teacher"
)

func TestDataStore_AddTeacher(t *testing.T) {
	expected := store.Teacher{
		Preferences: store.TeacherPreferences{
			TimeSlots: []store.TimeSlot{
				{
					Weekday:  time.Tuesday,
					Start:    timetype.NewUTCClock(20, 0, 0, 15),
					Duration: timetype.Duration(1*time.Hour + 30*time.Minute),
					Location: "room #108",
				},
			},
		},
		TeacherDetails: store.TeacherDetails{
			Name:    "Foo",
			Surname: "Barovich",
			Email:   "f.barovich@innopolis.university",
			Degree:  "Dr.",
			About:   "does it really matter?",
		},
	}
	srv := DataStore{TeacherRepository: &teacher.InterfaceMock{
		AddTeacherFunc: func(teacher store.TeacherDetails) (string, error) {
			assert.NotEmpty(t, teacher.ID)
			expected.ID = teacher.ID
			assert.Equal(t, expected.TeacherDetails, teacher)
			return expected.ID, nil
		},
		SetPreferencesFunc: func(teacherID string, pref store.TeacherPreferences) error {
			assert.Equal(t, expected.ID, teacherID)
			assert.Equal(t, expected.Preferences, pref)
			return nil
		},
	}}
	id, err := srv.AddTeacher(expected)
	require.NoError(t, err)
	assert.Equal(t, expected.ID, id)
}

func TestDataStore_AddAndRegisterUser(t *testing.T) {
	privs := []store.Privilege{store.PrivAddUsers, store.PrivEditUsers, store.PrivListUsers, store.PrivReadUsers}
	expected := store.User{
		Email:      "foo@bar.com",
		Privileges: []store.Privilege{store.PrivAddUsers, store.PrivEditUsers, store.PrivListUsers, store.PrivReadUsers},
	}
	expectedPwd := "some very strong password"

	// add user
	ur := &user.InterfaceMock{AddUserFunc: func(user store.User, pwd string, ignoreIfExists bool) (string, error) {
		assert.NotEmpty(t, user.ID)

		// to make structs be equal to user the single assert.Equal
		expected.ID = user.ID

		// as the order of privileges may be different
		assert.ElementsMatch(t, privs, user.Privileges)
		expected.Privileges = nil
		user.Privileges = nil

		assert.Equal(t, expected, user)

		err := bcrypt.CompareHashAndPassword([]byte(pwd), []byte(expectedPwd))
		require.NoError(t, err)
		assert.False(t, ignoreIfExists)
		return expected.ID, nil
	}}
	srv := DataStore{UserRepository: ur, BCryptCost: 4}

	id, err := srv.AddUser(expected, expectedPwd)
	require.NoError(t, err)
	assert.Equal(t, expected.ID, id)

	// register admin
	expected.Privileges = privs
	ur.AddUserFunc = func(user store.User, pwd string, ignoreIfExists bool) (string, error) {
		assert.NotEmpty(t, user.ID)

		// to make structs be equal to user the single assert.Equal
		expected.ID = user.ID

		// as the order of privileges may be different
		assert.ElementsMatch(t, privs, user.Privileges)
		user.Privileges = nil
		expected.Privileges = nil

		assert.Equal(t, expected, user)

		err := bcrypt.CompareHashAndPassword([]byte(pwd), []byte(expectedPwd))
		require.NoError(t, err)
		assert.True(t, ignoreIfExists)
		return expected.ID, nil
	}

	id, err = srv.RegisterAdmin(expected.Email, expectedPwd)
	require.NoError(t, err)
	assert.Equal(t, expected.ID, id)
}

func TestDataStore_PassThroughMethods(t *testing.T) {
	usr := store.User{
		ID:         "some awesome userID",
		Email:      "foo@bar.com",
		Privileges: []store.Privilege{store.PrivReadUsers, store.PrivListUsers},
	}
	tch := []store.Teacher{
		{
			TeacherDetails: store.TeacherDetails{
				ID:      "some awesome teacherID",
				Name:    "foo",
				Surname: "bar",
				Email:   "foo@bar.com",
				Degree:  "graduate",
				About:   "some details about teacher",
			},
			Preferences: store.TeacherPreferences{
				TimeSlots: []store.TimeSlot{
					{
						Weekday:  time.Monday,
						Start:    timetype.NewClock(20, 0, 0, 0, time.UTC),
						Duration: timetype.Duration(1*time.Hour + 30*time.Minute),
						Location: "room 108",
					},
					{
						Weekday:  time.Tuesday,
						Start:    timetype.NewClock(10, 0, 0, 0, time.UTC),
						Duration: timetype.Duration(1*time.Hour + 30*time.Minute),
						Location: "room 109",
					},
					{
						Weekday:  time.Friday,
						Start:    timetype.NewClock(15, 0, 0, 0, time.UTC),
						Duration: timetype.Duration(1*time.Hour + 30*time.Minute),
						Location: "room 102",
					},
				},
				Locations: []store.Location{"108", "102", "109"},
			},
		},
		{
			TeacherDetails: store.TeacherDetails{
				ID:      "some awesome second teacherID",
				Name:    "bar",
				Surname: "foo",
				Email:   "bar@foo.com",
				Degree:  "undergraduate",
				About:   "does it really matter?",
			},
		},
	}
	pref := store.TeacherPreferences{
		TimeSlots: []store.TimeSlot{
			{
				Weekday:  time.Sunday,
				Start:    timetype.NewClock(13, 0, 0, 0, time.UTC),
				Duration: timetype.Duration(1*time.Hour + 30*time.Minute),
				Location: "room 108",
			},
			{
				Weekday:  time.Wednesday,
				Start:    timetype.NewClock(19, 0, 0, 0, time.UTC),
				Duration: timetype.Duration(1*time.Hour + 30*time.Minute),
				Location: "room 109",
			},
			{
				Weekday:  time.Thursday,
				Start:    timetype.NewClock(11, 0, 0, 0, time.UTC),
				Duration: timetype.Duration(1*time.Hour + 30*time.Minute),
				Location: "room 102",
			},
		},
		Locations: []store.Location{"108", "102", "109"},
	}
	grps := []store.Group{
		{
			ID:        "some awesome group ID",
			Name:      "B20-01",
			StudyYear: store.StudyYear{ID: "some awesome sy ID", Name: "BS - Year 1 (Computer Engineering)"},
		},
		{
			ID:        "some awesome group ID 2",
			Name:      "B20-02",
			StudyYear: store.StudyYear{ID: "some awesome sy ID", Name: "BS - Year 1 (Computer Engineering)"},
		},
		{
			ID:        "some awesome group ID 3",
			Name:      "B19-03",
			StudyYear: store.StudyYear{ID: "some awesome sy ID 2", Name: "BS - Year 2"},
		},
	}

	srv := DataStore{UserRepository: &user.InterfaceMock{
		GetUserFunc: func(id string) (store.User, error) {
			assert.Equal(t, usr.ID, id)
			return usr, nil
		},
	}, TeacherRepository: &teacher.InterfaceMock{
		DeleteTeacherFunc: func(teacherID string) error {
			assert.Equal(t, tch[0].ID, teacherID)
			return nil
		},
		GetTeacherFullFunc: func(teacherID string) (store.Teacher, error) {
			assert.Equal(t, tch[0].ID, teacherID)
			return tch[0], nil
		},
		SetPreferencesFunc: func(teacherID string, p store.TeacherPreferences) error {
			assert.Equal(t, tch[0].ID, teacherID)
			assert.Equal(t, pref, p)
			return nil
		},
		ListTeachersFunc: func() ([]store.TeacherDetails, error) {
			return []store.TeacherDetails{tch[0].TeacherDetails, tch[1].TeacherDetails}, nil
		},
	}, GroupRepository: &uni.InterfaceMock{
		AddGroupFunc: func(g store.Group) (string, error) {
			assert.NotEmpty(t, g.ID)
			grps[0].ID = g.ID
			assert.Equal(t, grps[0].Name, g.Name)
			assert.Equal(t, grps[0].StudyYear.ID, g.StudyYear.ID)
			return g.ID, nil
		},
		DeleteGroupFunc: func(id string) error {
			assert.Equal(t, grps[1].ID, id)
			return nil
		},
		ListGroupsFunc: func() ([]store.Group, error) {
			return grps, nil
		},
	}, BCryptCost: 4}
	err := srv.DeleteTeacher(tch[0].ID)
	require.NoError(t, err)

	lst, err := srv.ListTeachers()
	require.NoError(t, err)
	assert.Equal(t, []store.TeacherDetails{tch[0].TeacherDetails, tch[1].TeacherDetails}, lst)

	tt, err := srv.GetTeacherFull(tch[0].ID)
	require.NoError(t, err)
	assert.Equal(t, tch[0], tt)

	err = srv.SetTeacherPreferences(tch[0].ID, pref)
	require.NoError(t, err)

	email, err := srv.GetUserEmail(usr.ID)
	require.NoError(t, err)
	assert.Equal(t, usr.Email, email)

	privs, err := srv.GetUserPrivs(usr.ID)
	require.NoError(t, err)
	assert.Equal(t, usr.Privileges, privs)

	id, err := srv.AddGroup(grps[0].Name, "some awesome sy ID")
	require.NoError(t, err)
	assert.Equal(t, grps[0].ID, id)

	g, err := srv.ListGroups()
	require.NoError(t, err)
	assert.ElementsMatch(t, grps, g)

	err = srv.DeleteGroup(grps[1].ID)
	require.NoError(t, err)
}

func TestDataStore_CheckUserCredentials(t *testing.T) {
	b, err := bcrypt.GenerateFromPassword([]byte("some very protected pwd"), 4)
	require.NoError(t, err)

	srv := DataStore{UserRepository: &user.InterfaceMock{
		GetPasswordHashFunc: func(email string) (string, error) {
			assert.Equal(t, "foo@bar.cc", email)
			return string(b), err
		},
	}}
	ok, err := srv.CheckUserCredentials("foo@bar.cc", "some very protected pwd")
	require.NoError(t, err)
	assert.True(t, ok)
}
