import React from "react";

import { LayerPopup } from "lib/popup";
import PopupContainer from "./PopupContainer";
// import DeleteNodePopup from "./DeleteNodePopup";
// import AddNodePopup from "./AddNodePopup";
// import EditNodePopup from "./EditNodePopup";
// import AddChannelPopup from "./AddChannelPopup";
// import EditChannelPopup from "./EditChannelPopup";
// import DeleteChannelPopup from "./DeleteChannelPopup";
// import AddUsersPopup from "./AddUsersPopup";
// import Add3rdPartyPopup from "./Add3rdPartyPopup";
// import EditUsersPopup from "./EditUsersPopup";
// import Edit3rdPartyPopup from "./Edit3rdPartyPopup";
// import DeleteUsersPopup from "./DeleteUsersPopup";
import SettingNodePopup from "./SettingNodePopup";
import SettingChannelPopup from "./SettingChannelPopup";
import DeletePopup from "./DeletePopup";
import AlertPopup from "./AlertPopup";
import SuccessPopup from "./SuccessPopup";
import UserPopup from "./UserPopup";

import Edit3rdPartyPopup from "./Edit3rdPartyPopup";
import EditUserPopup from "./EditUserPopup";
import AddUserPopup from "./AddUserPopup";
import ChangePasswordPopup from "./ChangePasswordPopup";
import TxListSearchPopup from "./TxListSearchPopup";
import SessionTimeOutPopup from "./SessionTimeOutPopup";

export default class Popup {
    static adduserPopup(props) {
        return LayerPopup.show(
            <PopupContainer>
                <AddUserPopup {...props} />
            </PopupContainer>
        );
    }

    static edit3rdpartyPopup(id, selectUser, source, type, className) {
        return LayerPopup.show(
            <PopupContainer>
                <Edit3rdPartyPopup
                    id={id}
                    selectUser={selectUser}
                    source={source}
                    type={type}
                    className={className}
                />
            </PopupContainer>
        );
    }

    static edituserPopup(id, source, type, selectUser, className) {
        return LayerPopup.show(
            <PopupContainer>
                <EditUserPopup
                    id={id}
                    selectUser={selectUser}
                    source={source}
                    type={type}
                    className={className}
                />
            </PopupContainer>
        );
    }

    //   static edit3rdpartyPopup(id, selectUser, callback) {
    //     return LayerPopup.show(
    //       <PopupContainer>
    //         <Edit3rdPartyPopup
    //           id={id}
    //           selectUser={selectUser}
    //           callback={callback}
    //         />
    //       </PopupContainer>
    //     );
    //   }
    static sessionTimeoutPopup(props) {
        return LayerPopup.show(
            <PopupContainer>
                <SessionTimeOutPopup {...props} />
            </PopupContainer>
        );
    }
    static deletePopup(props) {
        return LayerPopup.show(
            <PopupContainer>
                <DeletePopup {...props} />
            </PopupContainer>
        );
    }
    static settingNodePopup(props) {
        return LayerPopup.show(
            <PopupContainer>
                <SettingNodePopup {...props} />
            </PopupContainer>
        );
    }
    static settingChannelPopup(props) {
        return LayerPopup.show(
            <PopupContainer>
                <SettingChannelPopup {...props} />
            </PopupContainer>
        );
    }
    static successPopup(props) {
        return LayerPopup.show(
            <PopupContainer>
                <SuccessPopup {...props} />
            </PopupContainer>
        );
    }
    static alertPopup(props) {
        return LayerPopup.show(
            <PopupContainer>
                <AlertPopup {...props} />
            </PopupContainer>
        );
    }
    static txListSearchPopup(props) {
        return LayerPopup.show(
            <PopupContainer>
                <TxListSearchPopup {...props} />
            </PopupContainer>
        );
    }
    static userPopup(props) {
        return LayerPopup.show(
            <PopupContainer>
                <UserPopup {...props} />
            </PopupContainer>
        );
    }

    static changePasswordPopup(props) {
        return LayerPopup.show(
            <PopupContainer>
                <ChangePasswordPopup {...props} />
            </PopupContainer>
        );
    }
}
