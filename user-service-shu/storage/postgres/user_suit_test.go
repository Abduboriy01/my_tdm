package postgres

// import (
// //	"container/list"
// 	"testing"

// 	"github.com/my_tdm/user-service-shu/config"
// 	"github.com/my_tdm/user-service-shu/pkg/db"
// 	"github.com/my_tdm/user-service-shu/storage/repo"
// 	pb "github.com/my_tdm/user-service-shu/genproto"
// 	"github.com/stretchr/testify/suite"
// )

// type UserTestSuitRepo struct {
// 	suite.Suite
// 	CleanUpFunc func()
// 	Repository repo.UserStorageI
// }

// func (suite *UserTestSuitRepo) SetupSuite() {
// 	pgPool, cleanUp := db.ConnectToDB(config.Load())
// 	suite.Repository = NewUserRepo(pgPool)
// 	suite.CleanUpFunc = cleanUp
// }

// func (suite *UserTestSuitRepo) TestCreateCRUD(t *testing.T){
// 	user := pb.User{
// 		Id: "",
// 		FirstName: "Jamshid",
// 		LastName: "ikrom",
// 	}
// 	// create test
// 	createdUser, err := suite.Repository.CreateUser(&user)
// 	suite.Nil(err)
// 	suite.NotNil(createdUser)

// 	// get by id test
// 	getUser, err := suite.Repository.GetUserById(createdUser.Id)
// 	suite.Nil(err)
// 	suite.NotNil(getUser)

// 	getUser.FirstName = "Javohir"

// 	// update user test
// 	updatedUser, err := suite.Repository.Update(getUser)
// 	suite.Nil(err)
// 	suite.NotNil(updatedUser)

// 	// checikng if update is working
// 	getUpdatedUser, err := suite.Repository.GetUserById(updatedUser.Id)
// 	suite.Nil(err)
// 	suite.NotNil(getUser)
// 	suite.Equal(updatedUser.FirstName, getUpdatedUser.FirstName, "user should be equal after update")

// 	// list users test
// 	listUsers, _, err := suite.Repository.GetUserList(1, 10)
// 	suite.Nil(err)
// 	suite.NotNil(listUsers)

// 	err = suite.Repository.Delete(getUpdatedUser.Id)
// 	suite.Nil(err)
// }

// func (suite *UserTestSuitRepo) TearDownSuite() {
// 	suite.CleanUpFunc()
// }

// func TestUserRepositoryTestSuite(t *testing.T) {
// 	suite.Run(t, new(UserTestSuitRepo))
// }
