package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

type Student struct {
	gorm.Model
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Courses []Course `gorm:"many2many:student_courses;" json:"courses"`
}

type Course struct {
	gorm.Model
	Name     string    `json:"name"`
	Students []Student `gorm:"many2many:student_courses;" json:"students"`
}

type StudentCourse struct {
	gorm.Model
	StudentID uint
	CourseID  uint
}

// CRUD endpoints for students

func listStudents(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var students []Student
		if err := db.Preload("Courses").Find(&students).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, students)
	}
}

func getStudent(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
			return
		}
		var student Student
		if err := db.Preload("Courses").First(&student, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
			return
		}
		c.JSON(http.StatusOK, student)
	}
}

func createStudent(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var student Student
		if err := c.BindJSON(&student); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Create(&student).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, student)
	}

}

// updateStudent updates a student by ID
func updateStudent(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
			return
		}
		var student Student
		if err := db.First(&student, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
			return
		}
		if err := c.BindJSON(&student); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Save(&student).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, student)
	}
}

// deleteStudent deletes a student by ID
func deleteStudent(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
			return
		}
		if err := db.Delete(&Student{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Student deleted"})
	}
}

// CRUD endpoints for courses

func listCourses(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var courses []Course
		if err := db.Preload("Students").Find(&courses).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, courses)
	}
}

func getCourse(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
			return
		}
		var course Course
		if err := db.Preload("Students").First(&course, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
			return
		}
		c.JSON(http.StatusOK, course)
	}
}

func createCourse(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var course Course
		if err := c.BindJSON(&course); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Create(&course).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, course)
	}
}

func updateCourse(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
			return
		}
		var course Course
		if err := db.First(&course, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
			return
		}
		if err := c.BindJSON(&course); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := db.Save(&course).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, course)
	}
}

func deleteCourse(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
			return
		}
		if err := db.Delete(&Course{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Course deleted"})
	}
}

// endpoints for listing courses taken by a student and students taking a course

func listCoursesByStudent(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
			return
		}
		var student Student
		if err := db.Preload("Courses").First(&student, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
			return
		}
		c.JSON(http.StatusOK, student.Courses)
	}
}

func listStudentsByCourse(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
			return
		}
		var course Course
		if err := db.Preload("Students").First(&course, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
			return
		}
		c.JSON(http.StatusOK, course.Students)
	}
}

// updateCoursesByStudent updates the list of courses taken by a student
func updateCoursesByStudent(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		studentID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
			return
		}
		var student Student
		if err := db.Preload("Courses").First(&student, studentID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
			return
		}
		var courses []Course
		if err := c.BindJSON(&courses); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// delete all existing courses for the student
		if err := db.Model(&student).Association("Courses").Clear(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// add the new courses to the student's list
		for _, course := range courses {
			if err := db.Model(&student).Association("Courses").Append(&course); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{"message": "Courses updated"})
	}
}

func main() {
	dsn := "root:password@tcp(127.0.0.1:3306)/school_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(&Student{}, &Course{}, &StudentCourse{}); err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	// CRUD endpoints for students
	r.GET("/students", listStudents(db))
	r.GET("/students/:id", getStudent(db))
	r.POST("/students", createStudent(db))
	r.PUT("/students/:id", updateStudent(db))
	r.DELETE("/students/:id", deleteStudent(db))

	// CRUD endpoints for courses
	r.GET("/courses", listCourses(db))
	r.GET("/courses/:id", getCourse(db))
	r.POST("/courses", createCourse(db))
	r.PUT("/courses/:id", updateCourse(db))
	r.DELETE("/courses/:id", deleteCourse(db))

	// endpoints for listing courses taken by a student, listing students taking a particular course, and updating student courses
	r.GET("/students/:id/courses", listCoursesByStudent(db))
	r.GET("/courses/:id/students", listStudentsByCourse(db))
	r.PUT("/students/:id/courses", updateCoursesByStudent(db))
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
