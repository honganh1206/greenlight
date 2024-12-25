package main

import (
	"net/http"
	"testing"

	"greenlight.honganhpham.net/internal/assert"
)

func TestCreateMovieHandler(t *testing.T) {

	tl := newTestLogger(t)
	app := newTestApplication(t, tl)
	ts := newTestServer(t, app)

	defer ts.Close()

	tests := []struct {
		name           string
		inputJSON      string
		expectedStatus int
	}{
		{
			name: "Valid Input",
			inputJSON: `{
				"title": "Test Movie",
				"year": 2020,
				"runtime": "120 mins",
				"genres": ["drama", "action"]
			}`,
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Missing Title",
			inputJSON: `{
				"year": 2020,
				"runtime": "120 mins",
				"genres": ["drama", "action"]
			}`,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "Invalid Year",
			inputJSON: `{
				"title": "Test Movie",
				"year": 1800,
				"runtime": "120 mins",
				"genres": ["drama", "action"]
			}`,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "Too Many Genres",
			inputJSON: `{
				"title": "Test Movie",
				"year": 2020,
				"runtime": "120 mins",
				"genres": ["drama", "action", "comedy", "thriller", "horror", "documentary"]
			}`,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "Duplicate Genres",
			inputJSON: `{
				"title": "Test Movie",
				"year": 2020,
				"runtime": "120 mins",
				"genres": ["drama", "drama"]
			}`,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "Invalid JSON",
			inputJSON:      `{"title": "Test Movie"`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Empty Runtime",
			inputJSON: `{
            "title": "Test Movie",
            "year": 2020,
            "runtime": "",
            "genres": ["drama"]
        }`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Future Year",
			inputJSON: `{
            "title": "Test Movie",
            "year": 2525,
            "runtime": "120 mins",
            "genres": ["drama"]
        }`,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "Empty Genres Array",
			inputJSON: `{
            "title": "Test Movie",
            "year": 2020,
            "runtime": "120 mins",
            "genres": []
        }`,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "Invalid Runtime Format",
			inputJSON: `{
            "title": "Test Movie",
            "year": 2020,
            "runtime": "invalid",
            "genres": ["drama"]
        }`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, _ := ts.post(t, MovieV1, []byte(tt.inputJSON))
			assert.Equal(t, code, tt.expectedStatus)
		})
	}
}

func TestShowMovieHandler(t *testing.T) {
	tl := newTestLogger(t)
	app := newTestApplication(t, tl)
	ts := newTestServer(t, app)

	defer ts.Close()

	tests := []struct {
		name           string
		urlPath        string
		expectedStatus int
	}{
		{
			name:           "Valid ID",
			urlPath:        "/v1/movies/1",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid ID Format",
			urlPath:        "/v1/movies/abc",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Negative ID",
			urlPath:        "/v1/movies/-1",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Zero ID",
			urlPath:        "/v1/movies/0",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Very Large ID",
			urlPath:        "/v1/movies/999999999",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Invalid Path",
			urlPath:        "/v1/movies/1/invalid",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, _ := ts.get(t, tt.urlPath)
			assert.Equal(t, code, tt.expectedStatus)
		})
	}
}

func TestUpdateMovieHandler(t *testing.T) {
	tl := newTestLogger(t)
	app := newTestApplication(t, tl)
	ts := newTestServer(t, app)

	defer ts.Close()

	tests := []struct {
		name           string
		urlPath        string
		inputJSON      string
		expectedStatus int
	}{
		{
			name:    "Valid Update",
			urlPath: "/v1/movies/1",
			inputJSON: `{
                "title": "Updated Movie",
                "year": 2021,
                "runtime": "130 mins",
                "genres": ["drama", "sci-fi"]
            }`,
			expectedStatus: http.StatusOK,
		},
		{
			name:    "Invalid ID",
			urlPath: "/v1/movies/abc",
			inputJSON: `{
                "title": "Updated Movie",
                "year": 2021,
                "runtime": "130 mins",
                "genres": ["drama", "sci-fi"]
            }`,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:    "Invalid Input",
			urlPath: "/v1/movies/1",
			inputJSON: `{
                "year": 1800,
                "runtime": "130 mins",
                "genres": ["drama", "drama"]
            }`,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "Malformed JSON",
			urlPath:        "/v1/movies/1",
			inputJSON:      `{"title": "Bad JSON"`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Empty Request Body",
			urlPath:        "/v1/movies/1",
			inputJSON:      `{}`,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:    "Invalid Runtime Format",
			urlPath: "/v1/movies/1",
			inputJSON: `{
            "title": "Updated Movie",
            "year": 2021,
            "runtime": "invalid",
            "genres": ["drama"]
        }`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:    "Future Year",
			urlPath: "/v1/movies/1",
			inputJSON: `{
            "title": "Updated Movie",
            "year": 2525,
            "runtime": "120 mins",
            "genres": ["drama"]
        }`,
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, _ := ts.update(t, tt.urlPath, []byte(tt.inputJSON))
			assert.Equal(t, code, tt.expectedStatus)
		})
	}
}

func TestDeleteMovieHandler(t *testing.T) {
	tl := newTestLogger(t)
	app := newTestApplication(t, tl)

	ts := newTestServer(t, app)

	defer ts.Close()
	tests := []struct {
		name           string
		urlPath        string
		expectedStatus int
	}{
		{
			name:           "Valid Delete",
			urlPath:        "/v1/movies/1",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid ID Format",
			urlPath:        "/v1/movies/abc",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Zero ID",
			urlPath:        "/v1/movies/0",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Invalid Path",
			urlPath:        "/v1/movies/1/invalid",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			code, _, _ := ts.delete(t, tt.urlPath)
			assert.Equal(t, code, tt.expectedStatus)
		})
	}
}

// func TestMovieModel_Insert(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	defer db.Close()

// 	movieModel := MovieModel{DB: db}

// 	tests := []struct {
// 		name    string
// 		movie   *Movie
// 		mockErr error
// 		wantErr bool
// 	}{
// 		{
// 			name: "successful insert",
// 			movie: &Movie{
// 				Title:   "Test Movie",
// 				Year:    2023,
// 				Runtime: Runtime(120),
// 				Genres:  []string{"Action", "Drama"},
// 			},
// 			mockErr: nil,
// 			wantErr: false,
// 		},
// 		{
// 			name: "database error",
// 			movie: &Movie{
// 				Title:   "Test Movie",
// 				Year:    2023,
// 				Runtime: Runtime(120),
// 				Genres:  []string{"Action", "Drama"},
// 			},
// 			mockErr: sql.ErrConnDone,
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			rows := sqlmock.NewRows([]string{"id", "created_at", "version"}).
// 				AddRow(1, time.Now(), 1)

// 			mock.ExpectQuery(`INSERT INTO movies`).
// 				WithArgs(tt.movie.Title, tt.movie.Year, tt.movie.Runtime, pq.Array(tt.movie.Genres)).
// 				WillReturnRows(rows).
// 				WillReturnError(tt.mockErr)

// 			err := movieModel.Insert(tt.movie)
// 			if tt.wantErr {
// 				assert.Error(t, err)
// 			} else {
// 				assert.NoError(t, err)
// 				assert.NotZero(t, tt.movie.ID)
// 			}
// 		})
// 	}
// }

// func TestMovieModel_Get(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	defer db.Close()

// 	movieModel := MovieModel{DB: db}

// 	tests := []struct {
// 		name    string
// 		id      int64
// 		mockErr error
// 		want    *Movie
// 		wantErr error
// 	}{
// 		{
// 			name:    "valid id",
// 			id:      1,
// 			mockErr: nil,
// 			want: &Movie{
// 				ID:      1,
// 				Title:   "Test Movie",
// 				Year:    2023,
// 				Runtime: Runtime(120),
// 				Genres:  []string{"Action", "Drama"},
// 				Version: 1,
// 			},
// 			wantErr: nil,
// 		},
// 		{
// 			name:    "record not found",
// 			id:      999,
// 			mockErr: sql.ErrNoRows,
// 			want:    nil,
// 			wantErr: ErrRecordNotFound,
// 		},
// 		{
// 			name:    "invalid id",
// 			id:      0,
// 			mockErr: nil,
// 			want:    nil,
// 			wantErr: ErrRecordNotFound,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if tt.id > 0 && tt.mockErr != sql.ErrNoRows {
// 				rows := sqlmock.NewRows([]string{"id", "created_at", "title", "year", "runtime", "genres", "version"}).
// 					AddRow(tt.want.ID, time.Now(), tt.want.Title, tt.want.Year, tt.want.Runtime,
// 						pq.Array(tt.want.Genres), tt.want.Version)

// 				mock.ExpectQuery(`SELECT (.+) FROM movies`).
// 					WithArgs(tt.id).
// 					WillReturnRows(rows).
// 					WillReturnError(tt.mockErr)
// 			}

// 			got, err := movieModel.Get(tt.id)
// 			if tt.wantErr != nil {
// 				assert.ErrorIs(t, err, tt.wantErr)
// 			} else {
// 				assert.NoError(t, err)
// 				assert.Equal(t, tt.want.Title, got.Title)
// 				assert.Equal(t, tt.want.Year, got.Year)
// 				assert.Equal(t, tt.want.Runtime, got.Runtime)
// 				assert.Equal(t, tt.want.Genres, got.Genres)
// 			}
// 		})
// 	}
// }

// func TestMovieModel_Update(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	defer db.Close()

// 	movieModel := MovieModel{DB: db}

// 	tests := []struct {
// 		name    string
// 		movie   *Movie
// 		mockErr error
// 		wantErr error
// 	}{
// 		{
// 			name: "successful update",
// 			movie: &Movie{
// 				ID:      1,
// 				Title:   "Updated Movie",
// 				Year:    2024,
// 				Runtime: Runtime(130),
// 				Genres:  []string{"Action", "Sci-Fi"},
// 				Version: 1,
// 			},
// 			mockErr: nil,
// 			wantErr: nil,
// 		},
// 		{
// 			name: "edit conflict",
// 			movie: &Movie{
// 				ID:      1,
// 				Title:   "Updated Movie",
// 				Version: 1,
// 			},
// 			mockErr: sql.ErrNoRows,
// 			wantErr: ErrEditConflict,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if tt.mockErr != sql.ErrNoRows {
// 				rows := sqlmock.NewRows([]string{"version"}).AddRow(tt.movie.Version + 1)
// 				mock.ExpectQuery(`UPDATE movies`).
// 					WithArgs(tt.movie.Title, tt.movie.Year, tt.movie.Runtime,
// 						pq.Array(tt.movie.Genres), tt.movie.ID, tt.movie.Version).
// 					WillReturnRows(rows).
// 					WillReturnError(tt.mockErr)
// 			} else {
// 				mock.ExpectQuery(`UPDATE movies`).
// 					WithArgs(tt.movie.Title, tt.movie.Year, tt.movie.Runtime,
// 						pq.Array(tt.movie.Genres), tt.movie.ID, tt.movie.Version).
// 					WillReturnError(tt.mockErr)
// 			}

// 			err := movieModel.Update(tt.movie)
// 			if tt.wantErr != nil {
// 				assert.ErrorIs(t, err, tt.wantErr)
// 			} else {
// 				assert.NoError(t, err)
// 				assert.Equal(t, tt.movie.Version+1, tt.movie.Version)
// 			}
// 		})
// 	}
// }

// func TestMovieModel_Delete(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	defer db.Close()

// 	movieModel := MovieModel{DB: db}

// 	tests := []struct {
// 		name       string
// 		id         int64
// 		rowsAff    int64
// 		mockErr    error
// 		wantErr    error
// 	}{
// 		{
// 			name:    "successful delete",
// 			id:      1,
// 			rowsAff: 1,
// 			mockErr: nil,
// 			wantErr: nil,
// 		},
// 		{
// 			name:    "record not found",
// 			id:      999,
// 			rowsAff: 0,
// 			mockErr: nil,
// 			wantErr: ErrRecordNotFound,
// 		},
// 		{
// 			name:    "invalid id",
// 			id:      0,
// 			rowsAff: 0,
// 			mockErr: nil,
// 			wantErr: ErrRecordNotFound,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if tt.id > 0 {
// 				mock.ExpectExec(`DELETE FROM movies`).
// 					WithArgs(tt.id).
// 					WillReturnResult(sqlmock.NewResult(0, tt.rowsAff)).
// 					WillReturnError(tt.mockErr)
// 			}

// 			err := movieModel.Delete(tt.id)
// 			if tt.wantErr != nil {
// 				assert.ErrorIs(t, err, tt.wantErr)
// 			} else {
// 				assert.NoError(t, err)
// 			}
// 		})
// 	}
// }

// func TestMovieModel_GetAll(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	defer db.Close()

// 	movieModel := MovieModel{DB: db}

// 	tests := []struct {
// 		name    string
// 		title   string
// 		genres  []string
// 		filters Filters
// 		want    []*Movie
// 		mockErr error
// 		wantErr bool
// 	}{
// 		{
// 			name:   "successful retrieval",
// 			title:  "",
// 			genres: []string{},
// 			filters: Filters{
// 				Page:     1,
// 				PageSize: 20,
// 				Sort:     "id",
// 			},
// 			want: []*Movie{
// 				{
// 					ID:      1,
// 					Title:   "Movie 1",
// 					Year:    2023,
// 					Runtime: Runtime(120),
// 					Genres:  []string{"Action"},
// 					Version: 1,
// 				},
// 				{
// 					ID:      2,
// 					Title:   "Movie 2",
// 					Year:    2024,
// 					Runtime: Runtime(130),
// 					Genres:  []string{"Drama"},
// 					Version: 1,
// 				},
// 			},
// 			mockErr: nil,
// 			wantErr: false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			rows := sqlmock.NewRows([]string{
// 				"id", "created_at", "title", "year", "runtime", "genres", "version",
// 			})

// 			for _, movie := range tt.want {
// 				rows.AddRow(
// 					movie.ID,
// 					time.Now(),
// 					movie.Title,
// 					movie.Year,
// 					movie.Runtime,
// 					pq.Array(movie.Genres),
// 					movie.Version,
// 				)
// 			}

// 			mock.ExpectQuery(`SELECT (.+) FROM movies`).
// 				WillReturnRows(rows).
// 				WillReturnError(tt.mockErr)

// 			got, err := movieModel.GetAll(tt.title, tt.genres, tt.filters)
// 			if tt.wantErr {
// 				assert.Error(t, err)
// 			} else {
// 				assert.NoError(t, err)
// 				assert.Equal(t, len(tt.want), len(got))
// 				for i := range got {
// 					assert.Equal(t, tt.want[i].Title, got[i].Title)
// 					assert.Equal(t, tt.want[i].Year, got[i].Year)
// 					assert.Equal(t, tt.want[i].Runtime, got[i].Runtime)
// 					assert.Equal(t, tt.want[i].Genres, got[i].Genres)
// 				}
// 			}
// 		})
// 	}
// }
