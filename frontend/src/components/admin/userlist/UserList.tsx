import React, {FC, useContext, useEffect, useState} from "react";
import {FaUsers, FaEdit, FaUserPlus, FaTrashAlt} from "react-icons/fa";
import {MdDeleteForever, MdMoreHoriz} from "react-icons/md";
import {Button} from "@material-ui/core";
import {Loading, Popup} from "components";
import {UserListData} from "modules/users";
import {setRef} from "@material-ui/core/utils";
import Language from "../../common/Language";

interface Props {
    getUserList: () => any;
    deleteUser: (payload) => any;
    userList: UserListData;
    loading: boolean;
}

const UserList: FC<Props> = props => {
    const LanguageContext = useContext(Language());
    const {I18n} = LanguageContext;

    useEffect(() => {
        props.getUserList();
    }, []);

    const onClickAdd = () => {
        Popup.adduserPopup({className: "user"});
    }

    const onClick3rdEdit = (id, source, type, selectUser, className) => {
        Popup.edit3rdpartyPopup(id, selectUser, source, type, className);
    };

    const onClickEdit = (id, source, type, selectUser, className) => {
        Popup.edituserPopup(id, source, type, selectUser, className);
    };
    //
    const onClickDelete = userList => {
        Popup.deletePopup({
            className: "delete",
            name : userList.name,
            callbackFunc: () => {
                props.deleteUser({id : userList. id});
            }
        });
    };

    if (props.loading) {
        return <Loading/>;
    } else {
        return (
            <div className="content user-list">
                <div className="table-top">
                    {/*<i className="ic user-list-ic ic-type1">*/}
                    {/*<FaUsers/>*/}
                    {/*</i>*/}
                </div>
                <div className="table-box type-b">
                    <p className="title">
                        {I18n.users}
                        <Button className="add-user-btn" onClick={onClickAdd}>
                            {I18n.addCreateUsers}
                            <FaUserPlus/>
                        </Button>
                    </p>
                    <div className="table-content">
                        <table>
                            <thead>
                            <tr>
                                <th>{I18n.tableUser}</th>
                                <th>{I18n.tableUserType}</th>
                                <th>{I18n.tableEMail}</th>
                                <th>{I18n.tablePhone}</th>
                                <th>{I18n.tableChannel}</th>
                            </tr>
                            </thead>
                            <tbody>
                            {props.userList.data.map((el, key) => {
                                return (
                                    <tr key={key}>
                                        <th>{el.userId}</th>
                                        <th>{el.userType}</th>
                                        <th>{el.email}</th>
                                        <th>{el.phoneNumber}</th>
                                        <th>
                                            {el.channels.length} channel(s)
                                            <MdMoreHoriz className="more"/>
                                            <div className="channel-list">
                                                <ul>
                                                    {el.channels.map((el, key) => {
                                                        return <li key={key}>{el.name}</li>;
                                                    })}
                                                </ul>
                                            </div>
                                            {el.userType === "USER_ADMIN" && (
                                                <div className="btn-box">
                                                    <button className="edit-btn"
                                                            onClick={() => onClickEdit(el.id, "admin", "edit", el, "user-edit")}
                                                    >
                                                        <FaEdit/>
                                                    </button>
                                                </div>
                                            )}
                                            {el.userType === "USER_COMMON" && (
                                                <div className="btn-box">
                                                    <button className="edit-btn"
                                                            onClick={() => onClickEdit(el.id, "common", "edit", el, "user-edit")}
                                                    >
                                                        <FaEdit/>
                                                    </button>
                                                    <button className="delete-btn"
                                                            onClick={() => onClickDelete({
                                                                id:el.id,
                                                                name:el.userId
                                                            })}>
                                                        <FaTrashAlt />
                                                    </button>
                                                </div>
                                            )}
                                            {el.userType === "THIRD_PARTY" && (
                                                <div className="btn-box">
                                                    <button className="edit-btn"
                                                            onClick={() => onClick3rdEdit(el.id, "3rd", "edit", el, "third-user-edit")}
                                                    >
                                                        <FaEdit/>
                                                    </button>
                                                    <button className="delete-btn"
                                                            onClick={() => onClickDelete({
                                                                id:el.id,
                                                                name:el.userId
                                                            })}>
                                                        <FaTrashAlt />
                                                    </button>
                                                </div>
                                            )}
                                        </th>
                                    </tr>
                                );
                            })}
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        );
    }
};

export default UserList;

