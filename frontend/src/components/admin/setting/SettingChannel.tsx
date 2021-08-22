import React, { FC, useEffect, useContext } from "react";
import { FaEdit, FaTrashAlt } from "react-icons/fa";
import { MdLibraryAdd } from "react-icons/md";
import { ChannelsData } from "modules/channels";
import { Language, Popup } from "components";
import { makeSortList } from "utils/utils";

interface Props {
    loading: boolean;
    getChannels: () => any;
    channels: ChannelsData;
    deleteChannel: (payload) => any;
}

const SettingChannel: FC<Props> = props => {
    const LanguageContext = useContext(Language());
    const { I18n } = LanguageContext;
    useEffect(() => {
        props.getChannels();
    }, []);

    const onClickCreate = () => {
        Popup.settingChannelPopup({
            className: "setting-channel-popup",
            type: "create"
        });
    };

    const onClickEdit = channelInfo => {
        Popup.settingChannelPopup({
            className: "setting-channel-popup",
            type: "edit",
            channelInfo
        });
    };

    const onClickDelete = channelInfo => {
        Popup.deletePopup({
            className: "delete",
            name: channelInfo.name,
            callbackFunc: () => {
                props.deleteChannel({ id: channelInfo.channel });
            }
        });
    };

    return (
        <div className="content setting-channel setting">
            <p className="title">
                {I18n.settingChannel}
                <button className="create-btn" onClick={onClickCreate}>
                    <p>{I18n.addChannel}</p>
                    <MdLibraryAdd />
                </button>
            </p>
            <div className="channel-table-box">
                {makeSortList(props.channels.data).map((el, key) => {
                    return (
                        <div
                            key={key}
                            className="channel-table table-box type-b"
                        >
                            <div className="channel-top">
                                {el.name}
                                <div className="btn-box">
                                    <button
                                        className="edit-btn"
                                        onClick={() =>
                                            onClickEdit({
                                                channel: el.id,
                                                name: el.name
                                            })
                                        }
                                    >
                                        <FaEdit />
                                    </button>
                                    <button
                                        className="delete-btn"
                                        onClick={() =>
                                            onClickDelete({
                                                channel: el.id,
                                                name: el.name
                                            })
                                        }
                                    >
                                        <FaTrashAlt />
                                    </button>
                                </div>
                            </div>
                            <div className="table-content">
                                <table>
                                    <thead>
                                        <tr>
                                            <th>{I18n.nodeName}</th>
                                            <th>{I18n.nodeIP}</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {makeSortList(el.nodes).map(
                                            (el, key) => {
                                                return (
                                                    <tr key={key}>
                                                        <th>{el.name}</th>
                                                        <th>{el.ip}</th>
                                                    </tr>
                                                );
                                            }
                                        )}
                                    </tbody>
                                </table>
                            </div>
                        </div>
                    );
                })}
            </div>
        </div>
    );
};

export default SettingChannel;
