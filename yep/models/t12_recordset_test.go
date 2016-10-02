// Copyright 2016 NDP Systèmes. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package models

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateRecordSet(t *testing.T) {
	Convey("Test record creation", t, func() {
		env := NewEnvironment(1)
		Convey("Creating simple user John with no relations and checking ID", func() {
			userJohnData := FieldMap{
				"UserName": "John Smith",
				"Email":    "jsmith@example.com",
			}
			users := env.Pool("User").Call("Create", userJohnData).(RecordCollection)
			So(users.Len(), ShouldEqual, 1)
			So(users.Get("ID"), ShouldNotEqual, 0)
		})
		//Convey("Creating user Jane with related Profile", func() {
		//	userJaneProfileData := FieldMap{
		//		"Age":   23,
		//		"Money": 12345,
		//	}
		//	profile := env.Pool("Profile").Call("Create", userJaneProfileData).(RecordCollection)
		//	So(profile.Len(), ShouldEqual, 1)
		//	userJaneData := FieldMap{
		//		"UserName": "Jane Smith",
		//		"Email":    "jane.smith@example.com",
		//		"Profile":  profile,
		//	}
		//	userJane := env.Pool("User").Call("Create", &userJaneData).(RecordCollection)
		//	So(userJane.Len(), ShouldEqual, 1)
		//	So(userJane.Get("Profile").(RecordCollection).Get("ID"), ShouldEqual, profile.Get("ID"))
		//})
		//Convey("Creating a user Will Smith", func() {
		//	userWill := User_WithID{
		//		UserName: "Will Smith",
		//		Email:    "will.smith@example.com",
		//	}
		//	users := env.Create(&userWill)
		//	Convey("Created user ids should match struct's ID ", func() {
		//		So(len(users.Ids()), ShouldEqual, 1)
		//		So(users.ID(), ShouldEqual, userWill.ID)
		//	})
		//})
		//env.cr.Commit()
	})
}

//func TestSearchRecordSet(t *testing.T) {
//	Convey("Testing search through RecordSets", t, func() {
//		env := NewEnvironment(1)
//		Convey("Searching User Jane and getting struct through ReadOne", func() {
//			users := env.Pool("User").Filter("UserName", "=", "Jane Smith").Search()
//			So(len(users.Ids()), ShouldEqual, 1)
//			var userJane User_WithID
//			users.RelatedDepth(1).ReadOne(&userJane)
//			So(userJane.UserName, ShouldEqual, "Jane Smith")
//			So(userJane.Email, ShouldEqual, "jane.smith@example.com")
//			So(userJane.Profile.Age, ShouldEqual, 23)
//			So(userJane.Profile.Money, ShouldEqual, 12345)
//		})
//
//		Convey("Testing search all users and getting struct slice", func() {
//			usersAll := env.Pool("User").Search()
//			So(len(usersAll.Ids()), ShouldEqual, 3)
//			var userStructs []*User_PartialWithPosts
//			num := usersAll.ReadAll(&userStructs)
//			So(num, ShouldEqual, 3)
//			So(userStructs[0].Email, ShouldEqual, "jsmith@example.com")
//			So(userStructs[1].Email, ShouldEqual, "jane.smith@example.com")
//			So(userStructs[2].Email, ShouldEqual, "will.smith@example.com")
//		})
//		env.cr.Rollback()
//	})
//}
//
//func TestUpdateRecordSet(t *testing.T) {
//	Convey("Testing updates through RecordSets", t, func() {
//		env := NewEnvironment(1)
//		Convey("Simple update with params to user Jane", func() {
//			rsJane := env.Pool("User").Filter("UserName", "=", "Jane Smith").Search()
//			So(len(rsJane.Ids()), ShouldEqual, 1)
//			res := rsJane.Call("Write", FieldMap{"UserName": "Jane A. Smith"})
//			So(res, ShouldEqual, true)
//			var userJane User_WithID
//			rsJane.ReadOne(&userJane)
//			So(userJane.UserName, ShouldEqual, "Jane A. Smith")
//			So(userJane.Email, ShouldEqual, "jane.smith@example.com")
//		})
//		Convey("Simple update with struct", func() {
//			rsJane := env.Pool("User").Filter("UserName", "=", "Jane A. Smith").Search()
//			var userJohn User_WithID
//			rsJohn := env.Pool("User").Filter("UserName", "=", "John Smith")
//			rsJohn.ReadOne(&userJohn)
//			userJohn.Email = "jsmith2@example.com"
//			env.Sync(&userJohn)
//			var userJane2 User_WithID
//			rsJane.ReadOne(&userJane2)
//			So(userJane2.UserName, ShouldEqual, "Jane A. Smith")
//			So(userJane2.Email, ShouldEqual, "jane.smith@example.com")
//			var userJohn2 User_WithID
//			env.Pool("User").Filter("UserName", "=", "John Smith").ReadOne(&userJohn2)
//			So(userJohn2.UserName, ShouldEqual, "John Smith")
//			So(userJohn2.Email, ShouldEqual, "jsmith2@example.com")
//		})
//		env.cr.Commit()
//	})
//}
//
//func TestDeleteRecordSet(t *testing.T) {
//	env := NewEnvironment(1)
//	Convey("Delete user John Smith", t, func() {
//		users := env.Pool("User").Filter("UserName", "=", "John Smith")
//		num := users.Unlink()
//		Convey("Number of deleted record should be 1", func() {
//			So(num, ShouldEqual, 1)
//		})
//	})
//	env.cr.Rollback()
//}