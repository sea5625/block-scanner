import React, { FC, useEffect, useContext } from "react";
import { FaTrashAlt, FaEdit } from "react-icons/fa";
import { MdAddToQueue } from "react-icons/md";
import { NodesData } from "modules/nodes";
import { Language, Popup } from "components";
import { makeSortList } from "utils/utils";

interface Props {
    loading: boolean;
    allNodes: NodesData;
    getAllNodes: () => any;
    deleteNode: (payload) => any;
}

const SettingNode: FC<Props> = props => {
    const LanguageContext = useContext(Language());
    const { I18n } = LanguageContext;
    useEffect(() => {
        props.getAllNodes();
    }, []);

    const onClickCreate = () => {
        Popup.settingNodePopup({
            className: "setting-node-popup",
            type: "create"
        });
    };

    const onClickEdit = nodeInfo => {
        Popup.settingNodePopup({
            className: "setting-node-popup",
            type: "edit",
            nodeInfo
        });
    };

    const onClickDelete = nodeInfo => {
        Popup.deletePopup({
            className: "delete",
            name: nodeInfo.name,
            callbackFunc: () => {
                props.deleteNode({ id: nodeInfo.id });
            }
        });
    };

    return (
        <div className="content setting-node setting">
            <p className="title">
                {I18n.settingNode}
                <button className="create-btn" onClick={onClickCreate}>
                    <p>{I18n.addNode}</p>
                    <MdAddToQueue />
                </button>
            </p>
            <div className="node-table-box">
                <div className="node-table table-box type-b">
                    <div className="table-content">
                        <table>
                            <thead>
                                <tr>
                                    <th>{I18n.nodeName}</th>
                                    <th>{I18n.nodeIP}</th>
                                    <th></th>
                                </tr>
                            </thead>
                            <tbody>
                                {makeSortList(props.allNodes.data).map(
                                    (el, key) => {
                                        return (
                                            <tr key={key}>
                                                <th>{el.name}</th>
                                                <th>{el.ip}</th>
                                                <th>
                                                    <button
                                                        className="edit-btn"
                                                        onClick={() =>
                                                            onClickEdit({
                                                                name: el.name,
                                                                ip: el.ip,
                                                                id: el.id
                                                            })
                                                        }
                                                    >
                                                        <FaEdit />
                                                    </button>
                                                    <button
                                                        className="delete-btn"
                                                        onClick={() =>
                                                            onClickDelete({
                                                                name: el.name,
                                                                id: el.id
                                                            })
                                                        }
                                                    >
                                                        <FaTrashAlt />
                                                    </button>
                                                </th>
                                            </tr>
                                        );
                                    }
                                )}
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default SettingNode;
