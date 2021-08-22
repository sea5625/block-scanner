import { connect } from "react-redux";
import { Dispatch, Action } from "redux";
import { UserList } from "components";
import { actions as usersAction } from "modules/users";

const mapStateToProps = state => ({
  userList: state.users.userList,
  loading: state.users.userListLoading
});

const mapDispatchToProps = (dispatch: Dispatch<Action>) => ({
  getUserList: () => dispatch(usersAction.getUserList()),
  deleteUser:(payload) => dispatch(usersAction.deleteUser(payload))
});

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(UserList);
