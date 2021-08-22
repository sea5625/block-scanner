import React, { FC, useContext, useState, useEffect } from "react";
import { connect } from "react-redux";
import { Dispatch, Action } from "redux";
import { Button } from "@material-ui/core";
import { actions as usersActions } from "modules/users";
import { UserForm, UserAthority } from "components";

interface Props {
  title: string;
  loading: boolean;
}

const UserPopup: FC<Props> = props => {
  const [tab, setTab] = useState(0);

  return (
    <div className="popup-content clearfix">
      <p className="title">Add User</p>
      <div className="flex-box between">
        <UserForm />
        <UserAthority />
      </div>
      <div className="btn-box">
        <Button className="cancle-btn">Cancel</Button>
        <Button className="submit-btn">Submit</Button>
      </div>
    </div>
  );
};

const mapStateToProps = state => ({});
const mapDispatchToProps = (dispatch: Dispatch<Action>) => ({});

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(UserPopup);
