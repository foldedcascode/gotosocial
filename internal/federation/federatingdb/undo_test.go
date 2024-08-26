// GoToSocial
// Copyright (C) GoToSocial Authors admin@gotosocial.org
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package federatingdb_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/superseriousbusiness/activity/streams"
	"github.com/superseriousbusiness/gotosocial/testrig"
)

type UndoTestSuite struct {
	FederatingDBTestSuite
}

func (suite *UndoTestSuite) TestUndoLike() {

	// local_account_1 sent a follow request to remote_account_2;
	// remote_account_2 rejects the follow request
	favingAccount := suite.testAccounts["remote_account_1"]
	favedAccount := suite.testAccounts["local_account_1"]
	ctx := createTestContext(favingAccount, favedAccount)

	undo := streams.NewActivityStreamsUndo()

	// favedAccountURI := testrig.URLMustParse(favedAccount.URI)
	favingAccountURI := testrig.URLMustParse(favingAccount.URI)

	acceptActorProp := streams.NewActivityStreamsActorProperty()
	acceptActorProp.AppendIRI(favingAccountURI)
	undo.SetActivityStreamsActor(acceptActorProp)
	favToUndo := testrig.NewTestFaves()["remote_account_1_local_account_1_status_1"]
	op := streams.NewActivityStreamsObjectProperty()
	asFav, asErr := suite.tc.FaveToAS(ctx, favToUndo)
	suite.NoError(asErr)
	op.AppendActivityStreamsLike(asFav)
	undo.SetActivityStreamsObject(op)

	err := suite.federatingDB.Undo(ctx, undo)
	suite.NoError(err)

	// put the follow request in the database
	// fr := &gtsmodel.FollowRequest{
	// 	ID:              "01FJ1S8DX3STJJ6CEYPMZ1M0R3",
	// 	CreatedAt:       time.Now(),
	// 	UpdatedAt:       time.Now(),
	// 	URI:             uris.GenerateURIForFollow(followingAccount.Username, "01FJ1S8DX3STJJ6CEYPMZ1M0R3"),
	// 	AccountID:       followingAccount.ID,
	// 	TargetAccountID: followedAccount.ID,
	// }
	// err := suite.db.Put(ctx, fr)
	// suite.NoError(err)

	// asFollow, err := suite.tc.FollowToAS(ctx, suite.tc.FollowRequestToFollow(ctx, fr))
	// suite.NoError(err)

	// rejectingAccountURI := testrig.URLMustParse(followedAccount.URI)
	// requestingAccountURI := testrig.URLMustParse(followingAccount.URI)

	// // create a Reject
	// reject := streams.NewActivityStreamsReject()

	// // set the rejecting actor on it
	// acceptActorProp := streams.NewActivityStreamsActorProperty()
	// acceptActorProp.AppendIRI(rejectingAccountURI)
	// reject.SetActivityStreamsActor(acceptActorProp)

	// // Set the recreated follow as the 'object' property.
	// acceptObject := streams.NewActivityStreamsObjectProperty()
	// acceptObject.AppendActivityStreamsFollow(asFollow)
	// reject.SetActivityStreamsObject(acceptObject)

	// // Set the To of the reject as the originator of the follow
	// acceptTo := streams.NewActivityStreamsToProperty()
	// acceptTo.AppendIRI(requestingAccountURI)
	// reject.SetActivityStreamsTo(acceptTo)

	// // process the reject in the federating database
	// err = suite.federatingDB.Reject(ctx, reject)
	// suite.NoError(err)

	// // there should be nothing in the federator channel since nothing needs to be passed
	// _, ok := suite.getFederatorMsg(time.Second)
	// suite.False(ok)

	// // the follow request should not be in the database anymore -- it's been rejected
	// err = suite.db.GetByID(ctx, fr.ID, &gtsmodel.FollowRequest{})
	// suite.ErrorIs(err, db.ErrNoEntries)
}

func (suite *UndoTestSuite) TestUndoUnknownLike() {

	favingAccount := suite.testAccounts["remote_account_1"]
	favedAccount := suite.testAccounts["local_account_1"]
	ctx := createTestContext(favingAccount, favedAccount)

	undo := streams.NewActivityStreamsUndo()

	favingAccountURI := testrig.URLMustParse(favingAccount.URI)

	acceptActorProp := streams.NewActivityStreamsActorProperty()
	acceptActorProp.AppendIRI(favingAccountURI)
	undo.SetActivityStreamsActor(acceptActorProp)

	op := streams.NewActivityStreamsObjectProperty()

	apLike2 := streams.NewActivityStreamsLike()
	actorProp2 := streams.NewActivityStreamsActorProperty()
	actorProp2.AppendIRI(testrig.URLMustParse("http://fossbros-anonymous.io/users/foss_satan"))
	apLike2.SetActivityStreamsActor(actorProp2)
	idProp := streams.NewJSONLDIdProperty()
	idProp.Set(
		testrig.URLMustParse("http://fossbros-anonymous.io/users/foss_satan/liked/01F8Q0486ACGGWKG02A7DS1Q28"),
	)
	apLike2.SetJSONLDId(idProp)
	objectProp2 := streams.NewActivityStreamsObjectProperty()
	objectProp2.AppendIRI(
		testrig.URLMustParse("http://localhost:8080/users/the_mighty_zork/statuses/NOTAREALID"),
	)
	apLike2.SetActivityStreamsObject(objectProp2)
	toProp2 := streams.NewActivityStreamsToProperty()
	toProp2.AppendIRI(
		testrig.URLMustParse("http://localhost:8080/users/the_mighty_zork"),
	)
	apLike2.SetActivityStreamsTo(toProp2)
	op.AppendActivityStreamsLike(apLike2)

	undo.SetActivityStreamsObject(op)

	err := suite.federatingDB.Undo(ctx, undo)
	suite.NoError(err)
}

func TestUndoTestSuite(t *testing.T) {
	suite.Run(t, &UndoTestSuite{})
}
