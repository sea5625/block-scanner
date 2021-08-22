import React, { Fragment } from "react";
import { connect } from "react-redux";
import { Dispatch, Action } from "redux";
import { actions as authActions } from "modules/auth";

const Loading = props => {
  return (
    <div className="loading">
      <button onClick={props.logout}>logout</button>
    </div>
  );
};

const mapStateToProps = state => ({
  user: state.users.user,
  loading: state.users.loading
});
const mapDispatchToProps = (dispatch: Dispatch<Action>) => ({
  logout: () => dispatch(authActions.logout())
});

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(Loading);
