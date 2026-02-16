package filebased

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"a21hc3NpZ25tZW50/model"

	"go.etcd.io/bbolt"
)

type Data struct {
	DB *bbolt.DB
}

func InitDB() (*Data, error) {
	dbPath := os.Getenv("APP_DB_PATH")
	if dbPath == "" {
		dbPath = "file.db"
	}

	db, err := bbolt.Open(dbPath, 0600, &bbolt.Options{Timeout: 2 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Tasks"))
		if err != nil {
			return fmt.Errorf("create tasks bucket: %v", err)
		}

		categoriesBucket, err := tx.CreateBucketIfNotExists([]byte("Categories"))
		if err != nil {
			return fmt.Errorf("create categories bucket: %v", err)
		}

		// Add default categories if bucket is empty
		catCount := 0
		categoriesBucket.ForEach(func(_, _ []byte) error {
			catCount++
			return nil
		})

		if catCount == 0 {
			// No default system categories needed anymore
			// Each user will get their own categories when they register
			fmt.Println("DEBUG: No default system categories created, users will get individual categories")
		}

		_, err = tx.CreateBucketIfNotExists([]byte("Users"))
		if err != nil {
			return fmt.Errorf("create users bucket: %v", err)
		}

		_, err = tx.CreateBucketIfNotExists([]byte("Sessions"))
		if err != nil {
			return fmt.Errorf("create sessions bucket: %v", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &Data{DB: db}, nil
}

func (data *Data) StoreTask(task model.Task) error {
	// Check if we need to generate an ID
	if task.ID <= 0 {
		return data.DB.Update(func(tx *bbolt.Tx) error {
			b := tx.Bucket([]byte("Tasks"))

			// Find the highest existing ID
			maxID := 0
			err := b.ForEach(func(k, v []byte) error {
				id := btoi(k)
				if id > maxID {
					maxID = id
				}
				return nil
			})
			if err != nil {
				return err
			}

			// Assign the next ID
			task.ID = maxID + 1

			// Marshal and store
			taskJSON, err := json.Marshal(task)
			if err != nil {
				return err
			}
			return b.Put([]byte(fmt.Sprintf("%d", task.ID)), taskJSON)
		})
	}

	// If already has an ID
	taskJSON, err := json.Marshal(task)
	if err != nil {
		return err
	}
	return data.DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Tasks"))
		return b.Put([]byte(fmt.Sprintf("%d", task.ID)), taskJSON)
	})
}

func (data *Data) StoreCategory(category model.Category) error {
	// Check if we need to generate an ID
	if category.ID <= 0 {
		return data.DB.Update(func(tx *bbolt.Tx) error {
			b := tx.Bucket([]byte("Categories"))

			// Find the highest existing ID
			maxID := 0
			err := b.ForEach(func(k, v []byte) error {
				id := btoi(k)
				if id > maxID {
					maxID = id
				}
				return nil
			})
			if err != nil {
				return err
			}

			// Assign the next ID
			category.ID = maxID + 1

			// Marshal and store
			categoryJSON, err := json.Marshal(category)
			if err != nil {
				return err
			}
			return b.Put([]byte(fmt.Sprintf("%d", category.ID)), categoryJSON)
		})
	}

	// If already has an ID
	categoryJSON, err := json.Marshal(category)
	if err != nil {
		return err
	}
	return data.DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Categories"))
		return b.Put([]byte(fmt.Sprintf("%d", category.ID)), categoryJSON)
	})
}

func (data *Data) UpdateTask(id int, task model.Task) error {
	return data.StoreTask(task) // Reuse StoreTask as it will replace the existing entry
}

func (data *Data) UpdateCategory(id int, category model.Category) error {
	return data.StoreCategory(category) // Reuse StoreCategory as it will replace the existing entry
}

func (data *Data) DeleteTask(id int) error {
	return data.DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Tasks"))
		return b.Delete([]byte(fmt.Sprintf("%d", id)))
	})
}

func (data *Data) DeleteCategory(id int) error {
	return data.DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Categories"))
		return b.Delete([]byte(fmt.Sprintf("%d", id)))
	})
}

func (data *Data) GetTaskByID(id int) (*model.Task, error) {
	var task model.Task
	err := data.DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Tasks"))
		v := b.Get([]byte(fmt.Sprintf("%d", id)))
		if v == nil {
			return fmt.Errorf("record not found")
		}
		return json.Unmarshal(v, &task)
	})
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (data *Data) GetCategoryByID(id int) (*model.Category, error) {
	var category model.Category
	err := data.DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Categories"))
		v := b.Get([]byte(fmt.Sprintf("%d", id)))
		if v == nil {
			return fmt.Errorf("record not found")
		}
		return json.Unmarshal(v, &category)
	})
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (data *Data) GetTasks() ([]model.Task, error) {
	var tasks []model.Task
	err := data.DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Tasks"))
		return b.ForEach(func(k, v []byte) error {
			var task model.Task
			if err := json.Unmarshal(v, &task); err != nil {
				log.Println("Error unmarshaling task:", err)
				return nil // Continue despite error
			}
			tasks = append(tasks, task)
			return nil
		})
	})
	if err != nil {
		return nil, fmt.Errorf("error fetching tasks: %v", err)
	}
	return tasks, nil
}

func (data *Data) GetTasksByUserID(userID int) ([]model.Task, error) {
	var tasks []model.Task
	err := data.DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Tasks"))
		return b.ForEach(func(k, v []byte) error {
			var task model.Task
			if err := json.Unmarshal(v, &task); err != nil {
				log.Println("Error unmarshaling task:", err)
				return nil // Continue despite error
			}
			if task.UserID == userID {
				tasks = append(tasks, task)
			}
			return nil
		})
	})
	if err != nil {
		return nil, fmt.Errorf("error fetching tasks: %v", err)
	}
	return tasks, nil
}

func (data *Data) GetCategories() ([]model.Category, error) {
	var categories []model.Category
	err := data.DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Categories"))
		return b.ForEach(func(k, v []byte) error {
			var category model.Category
			if err := json.Unmarshal(v, &category); err != nil {
				log.Println("Error unmarshaling category:", err)
				return nil // Continue despite error
			}
			categories = append(categories, category)
			return nil
		})
	})
	if err != nil {
		return nil, fmt.Errorf("error fetching categories: %v", err)
	}
	return categories, nil
}

// GetCategoriesByUserID returns categories for specific user only
func (data *Data) GetCategoriesByUserID(userID int) ([]model.Category, error) {
	var categories []model.Category
	err := data.DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Categories"))
		return b.ForEach(func(k, v []byte) error {
			var category model.Category
			if err := json.Unmarshal(v, &category); err != nil {
				log.Println("Error unmarshaling category:", err)
				return nil // Continue despite error
			}
			// Only include categories that belong to this specific user
			if category.UserID == userID {
				categories = append(categories, category)
			}
			return nil
		})
	})
	if err != nil {
		return nil, fmt.Errorf("error fetching categories for user %d: %v", userID, err)
	}
	return categories, nil
}

func (data *Data) Reset() error {
	return data.DB.Update(func(tx *bbolt.Tx) error {
		if err := tx.DeleteBucket([]byte("Tasks")); err != nil {
			return err
		}
		if err := tx.DeleteBucket([]byte("Categories")); err != nil {
			return err
		}
		if err := tx.DeleteBucket([]byte("Users")); err != nil {
			return err
		}

		return nil
	})
}

func (data *Data) CloseDB() error {
	return data.DB.Close()
}

func (data *Data) GetTaskListByCategory(categoryID int) ([]model.TaskCategory, error) {
	var taskCategories []model.TaskCategory
	category, err := data.GetCategoryByID(categoryID)
	if err != nil {
		return nil, fmt.Errorf("error fetching category: %v", err)
	}

	err = data.DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Tasks"))
		if b == nil {
			return fmt.Errorf("tasks bucket not found")
		}
		return b.ForEach(func(k, v []byte) error {
			var task model.Task
			if err := json.Unmarshal(v, &task); err != nil {
				log.Printf("Error unmarshaling task: %v", err)
				return nil // Continue processing next item in case of error
			}
			if task.CategoryID == categoryID {
				taskCategories = append(taskCategories, model.TaskCategory{
					ID:       task.ID,
					Title:    task.Title,
					Category: category.Name,
				})
			}
			return nil
		})
	})
	if err != nil {
		return nil, fmt.Errorf("error fetching tasks for category %d: %v", categoryID, err)
	}
	if len(taskCategories) == 0 {
		return nil, fmt.Errorf("no tasks found for category ID: %d", categoryID)
	}
	return taskCategories, nil
}

// GetTaskListByCategoryAndUser returns tasks filtered by both category and user
func (data *Data) GetTaskListByCategoryAndUser(categoryID, userID int) ([]model.TaskCategory, error) {
	var taskCategories []model.TaskCategory
	category, err := data.GetCategoryByID(categoryID)
	if err != nil {
		return nil, fmt.Errorf("error fetching category: %v", err)
	}

	err = data.DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Tasks"))
		if b == nil {
			return fmt.Errorf("tasks bucket not found")
		}
		return b.ForEach(func(k, v []byte) error {
			var task model.Task
			if err := json.Unmarshal(v, &task); err != nil {
				log.Printf("Error unmarshaling task: %v", err)
				return nil // Continue processing next item in case of error
			}
			// Filter by both category and user
			if task.CategoryID == categoryID && task.UserID == userID {
				taskCategories = append(taskCategories, model.TaskCategory{
					ID:       task.ID,
					Title:    task.Title,
					Category: category.Name,
				})
			}
			return nil
		})
	})
	if err != nil {
		return nil, fmt.Errorf("error fetching tasks for category %d and user %d: %v", categoryID, userID, err)
	}
	return taskCategories, nil
}

func (data *Data) GetUserByEmail(email string) (model.User, error) {
	var user model.User
	found := false // Flag to check if the user is found

	err := data.DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Users"))
		if b == nil {
			return fmt.Errorf("users bucket not found")
		}

		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var u model.User
			if err := json.Unmarshal(v, &u); err != nil {
				continue // Skip on unmarshal error
			}
			if u.Email == email {
				user = u
				found = true
				break // Stop the loop once the user is found
			}
		}
		return nil // Return nil error from the View transaction
	})

	if err != nil {
		return model.User{}, err // Return the error if the transaction failed
	}
	if !found {
		return model.User{}, nil // Return an empty User struct and nil error if not found
	}
	return user, nil // Return the found user and nil error
}

func (data *Data) CreateUser(user model.User) (model.User, error) {
	err := data.DB.Update(func(tx *bbolt.Tx) error {
		usersBucket := tx.Bucket([]byte("Users"))
		if usersBucket == nil {
			return fmt.Errorf("users bucket not found")
		}

		// Cek email sudah ada atau belum
		emailExists := false
		_ = usersBucket.ForEach(func(k, v []byte) error {
			var u model.User
			if err := json.Unmarshal(v, &u); err == nil {
				if u.Email == user.Email {
					emailExists = true
				}
			}
			return nil
		})
		if emailExists {
			return fmt.Errorf("email already exists")
		}

		// Find the highest existing ID
		maxID := 0
		err := usersBucket.ForEach(func(k, v []byte) error {
			id := btoi(k)
			if id > maxID {
				maxID = id
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("error reading user IDs: %v", err)
		}

		// Set the new user's ID to maxID + 1
		newUserID := maxID + 1
		user.ID = newUserID // Assuming User.ID is of type int

		userJSON, err := json.Marshal(user)
		if err != nil {
			return fmt.Errorf("error marshaling user: %v", err)
		}

		// Store the new user with the new ID (pakai string angka sebagai key)
		return usersBucket.Put([]byte(strconv.Itoa(newUserID)), userJSON)
	})
	if err != nil {
		return model.User{}, err
	}

	// Create default categories for the new user
	err = data.CreateDefaultCategoriesForUser(user.ID)
	if err != nil {
		// Log error but don't fail user creation
		fmt.Printf("Warning: Failed to create default categories for user %d: %v\n", user.ID, err)
	}

	return user, nil
}

// itob converts an integer to a byte slice
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

// btoi converts a byte slice to an integer
func btoi(b []byte) int {
	// Convert string representation to integer
	s := string(b)
	i, err := strconv.Atoi(s)
	if err != nil {
		// Return 0 if conversion fails
		return 0
	}
	return i
}

func (data *Data) GetUserTaskCategory() ([]model.UserTaskCategory, error) {
	var results []model.UserTaskCategory

	// Debug counts
	var userCount, taskCount, categoryCount, matchedTaskCount int

	fmt.Println("DEBUG - Starting GetUserTaskCategory method...")

	err := data.DB.View(func(tx *bbolt.Tx) error {
		usersBucket := tx.Bucket([]byte("Users"))
		tasksBucket := tx.Bucket([]byte("Tasks"))
		categoriesBucket := tx.Bucket([]byte("Categories"))

		if usersBucket == nil || tasksBucket == nil || categoriesBucket == nil {
			fmt.Println("ERROR - One or more required buckets do not exist")
			return fmt.Errorf("one or more required buckets do not exist")
		}

		// Count users
		usersBucket.ForEach(func(k, v []byte) error {
			userCount++
			return nil
		})

		// Count tasks
		tasksBucket.ForEach(func(k, v []byte) error {
			taskCount++
			return nil
		})

		// Count categories
		categoriesBucket.ForEach(func(k, v []byte) error {
			categoryCount++
			return nil
		})

		fmt.Printf("DEBUG - Database counts: Users=%d, Tasks=%d, Categories=%d\n", userCount, taskCount, categoryCount)

		// List all users IDs first
		fmt.Println("DEBUG - Listing all user IDs in database:")
		usersBucket.ForEach(func(k, v []byte) error {
			var user model.User
			if err := json.Unmarshal(v, &user); err == nil {
				fmt.Printf("  - User found: ID=%d, Email=%s\n", user.ID, user.Email)
			}
			return nil
		})

		// List all tasks with their UserIDs
		fmt.Println("DEBUG - Listing all tasks with their UserIDs:")
		tasksBucket.ForEach(func(k, v []byte) error {
			var task model.Task
			if err := json.Unmarshal(v, &task); err == nil {
				fmt.Printf("  - Task found: ID=%d, Title=%s, UserID=%d, CategoryID=%d\n",
					task.ID, task.Title, task.UserID, task.CategoryID)
			}
			return nil
		})

		return usersBucket.ForEach(func(_, userValue []byte) error {
			var user model.User
			if err := json.Unmarshal(userValue, &user); err != nil {
				fmt.Printf("ERROR - Failed to unmarshal user data: %v\n", err)
				return nil // skip badly formatted user records
			}

			fmt.Printf("DEBUG - Processing user: ID=%d, Email=%s\n", user.ID, user.Email)
			userTaskCount := 0

			// Now fetch tasks for the user
			err := tasksBucket.ForEach(func(_, taskValue []byte) error {
				var task model.Task
				if err := json.Unmarshal(taskValue, &task); err != nil {
					fmt.Printf("ERROR - Failed to unmarshal task data: %v\n", err)
					return nil // skip badly formatted task records
				}

				fmt.Printf("DEBUG - Checking task: ID=%d, Title=%s, UserID=%d, Against User.ID=%d\n",
					task.ID, task.Title, task.UserID, user.ID)

				if task.UserID == user.ID { // Check if the task belongs to the user
					userTaskCount++
					matchedTaskCount++

					categoryName := "Unknown"
					catValue := categoriesBucket.Get([]byte(fmt.Sprintf("%d", task.CategoryID)))
					if catValue != nil {
						var category model.Category
						if err := json.Unmarshal(catValue, &category); err != nil {
							fmt.Printf("ERROR - Failed to unmarshal category data for ID %d: %v\n", task.CategoryID, err)
						} else {
							categoryName = category.Name
							fmt.Printf("DEBUG - Found category: ID=%d, Name=%s\n", task.CategoryID, categoryName)
						}
					} else {
						fmt.Printf("WARNING - No category found with ID: %d\n", task.CategoryID)
					}

					result := model.UserTaskCategory{
						ID:       int(user.ID),
						Fullname: user.Fullname,
						Email:    user.Email,
						Task:     task.Title,
						Deadline: task.Deadline,
						Priority: task.Priority,
						Status:   task.Status,
						Category: categoryName,
					}
					results = append(results, result)
					fmt.Printf("DEBUG - Added task to results: Task=%s, UserID=%d, Category=%s\n",
						task.Title, user.ID, categoryName)
				}
				return nil
			})

			fmt.Printf("DEBUG - User ID=%d has %d tasks\n", user.ID, userTaskCount)
			return err
		})
	})

	if err != nil {
		fmt.Printf("ERROR in GetUserTaskCategory: %v\n", err)
		return nil, err
	}

	fmt.Printf("DEBUG - GetUserTaskCategory finished. Found %d users, %d total tasks, and %d tasks matched to users. Returning %d results.\n",
		userCount, taskCount, matchedTaskCount, len(results))

	return results, nil
}

func (data *Data) AddSession(session model.Session) error {
	sessionJSON, err := json.Marshal(session)
	if err != nil {
		return err
	}
	return data.DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Sessions"))
		return b.Put([]byte(session.Token), sessionJSON)
	})
}

func (data *Data) DeleteSession(token string) error {
	return data.DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Sessions"))
		return b.Delete([]byte(token))
	})
}

func (data *Data) UpdateSession(session model.Session) error {
	return data.AddSession(session) // Reuse AddSession as it will overwrite the existing entry
}

func (data *Data) SessionByToken(token string) (model.Session, error) {
	var session model.Session
	err := data.DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Sessions"))
		sessionData := b.Get([]byte(token))
		if sessionData == nil {
			return fmt.Errorf("session not found")
		}
		return json.Unmarshal(sessionData, &session)
	})
	if err != nil {
		return model.Session{}, err
	}
	return session, nil
}

func (data *Data) TokenExpired(session model.Session) bool {
	return session.Expiry.Before(time.Now())
}

func (data *Data) TokenValidity(token string) (model.Session, error) {
	session, err := data.SessionByToken(token)
	if err != nil {
		return model.Session{}, err
	}

	if data.TokenExpired(session) {
		err := data.DeleteSession(token)
		if err != nil {
			return model.Session{}, err
		}
		return model.Session{}, fmt.Errorf("session expired")
	}

	return session, nil
}

func (data *Data) GetFirstSession() (model.Session, error) {
	var session model.Session
	found := false // Flag to check if at least one session found

	err := data.DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Sessions"))
		if b == nil {
			return fmt.Errorf("sessions bucket not found")
		}

		// Use Cursor to iterate through the bucket
		c := b.Cursor()
		k, v := c.First() // Retrieve the first session
		if k != nil {
			err := json.Unmarshal(v, &session)
			if err != nil {
				return err // Return unmarshaling error
			}
			found = true // Set found true as we have at least one session
		}
		return nil
	})

	if err != nil {
		return model.Session{}, err // Return error encountered during the View transaction
	}

	if !found {
		return model.Session{}, fmt.Errorf("no sessions found") // No session was found
	}

	return session, nil // Return the first session found
}

func (data *Data) SessionAvailEmail(email string) (model.Session, error) {
	var session model.Session
	found := false // Flag to check if at least one session matches the email

	err := data.DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Sessions"))
		if b == nil {
			return fmt.Errorf("sessions bucket not found")
		}

		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var s model.Session
			if err := json.Unmarshal(v, &s); err != nil {
				continue // Skip badly formatted session records
			}
			if s.Email == email {
				session = s
				found = true
				break // Stop the iteration as we found the session
			}
		}
		return nil
	})

	if err != nil {
		return model.Session{}, err // Return error encountered during the View transaction
	}

	if !found {
		return model.Session{}, fmt.Errorf("no session available for email: %s", email) // No session was found for the given email
	}

	return session, nil // Return the found session
}

func (data *Data) SessionAvailToken(token string) (model.Session, error) {
	var session model.Session
	err := data.DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Sessions"))
		if b == nil {
			return fmt.Errorf("sessions bucket not found")
		}

		sessionData := b.Get([]byte(token))
		if sessionData == nil {
			return fmt.Errorf("no session available for token: %s", token) // No session was found for the given token
		}

		return json.Unmarshal(sessionData, &session)
	})

	if err != nil {
		return model.Session{}, err // Return error encountered during the View transaction or unmarshaling
	}

	return session, nil // Return the found session
}

// GetUsers retrieves all users from the database
func (data *Data) GetUsers() ([]model.User, error) {
	users := []model.User{}

	err := data.DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Users"))
		if b == nil {
			return fmt.Errorf("users bucket not found")
		}

		// Iterate over all users in the bucket
		return b.ForEach(func(k, v []byte) error {
			var user model.User
			if err := json.Unmarshal(v, &user); err != nil {
				// Log the error but continue with the next user
				fmt.Printf("Error unmarshaling user data: %v\n", err)
				return nil
			}
			// Add the user to the slice
			users = append(users, user)
			return nil
		})
	})

	if err != nil {
		return nil, err
	}

	fmt.Printf("DEBUG - Retrieved %d users from database\n", len(users))
	return users, nil
}

// CreateDefaultCategoriesForUser creates default categories for a new user
func (data *Data) CreateDefaultCategoriesForUser(userID int) error {
	fmt.Printf("DEBUG: CreateDefaultCategoriesForUser called for userID: %d\n", userID)
	return data.DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Categories"))
		if b == nil {
			return fmt.Errorf("categories bucket not found")
		}

		// First, check if user already has ANY categories (including checking each one individually)
		userCategories := []model.Category{}
		err := b.ForEach(func(k, v []byte) error {
			var category model.Category
			if err := json.Unmarshal(v, &category); err == nil {
				if category.UserID == userID {
					userCategories = append(userCategories, category)
					fmt.Printf("DEBUG: Found existing category for user %d: %s (ID: %d)\n", userID, category.Name, category.ID)
				}
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("error checking existing categories: %v", err)
		}

		// If user already has categories, don't create more
		if len(userCategories) > 0 {
			fmt.Printf("DEBUG: User %d already has %d categories, skipping creation\n", userID, len(userCategories))
			return nil
		}

		fmt.Printf("DEBUG: User %d has no categories, creating initial categories\n", userID)

		// Get the highest existing category ID to avoid conflicts between users
		maxID := 0
		err = b.ForEach(func(k, v []byte) error {
			id := btoi(k)
			if id > maxID {
				maxID = id
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("error reading category IDs: %v", err)
		}

		fmt.Printf("DEBUG: Current maxID in database: %d\n", maxID)

		initialCategories := []string{"Work", "Personal", "Study", "Health", "Home"}

		for i, categoryName := range initialCategories {
			categoryID := maxID + i + 1 // Globally unique ID
			category := model.Category{
				ID:     categoryID,
				Name:   categoryName,
				UserID: userID,
			}

			fmt.Printf("DEBUG: Creating category: %s (ID: %d) for user %d\n", category.Name, category.ID, userID)
			catJSON, err := json.Marshal(category)
			if err != nil {
				return fmt.Errorf("error marshaling category: %v", err)
			}
			err = b.Put([]byte(fmt.Sprintf("%d", category.ID)), catJSON)
			if err != nil {
				return fmt.Errorf("error storing category: %v", err)
			}
		}

		fmt.Printf("DEBUG: Successfully created %d categories for user %d\n", len(initialCategories), userID)
		return nil
	})
}
