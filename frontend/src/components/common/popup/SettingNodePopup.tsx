import React, { FC, useState, useEffect, useContext } from "react";
import { connect } from "react-redux";
import { Dispatch, Action } from "redux";
import { Button } from "@material-ui/core";
import { actions as nodesActions } from "modules/nodes";
import { Language } from "components";
import { NodesData } from "modules/nodes";
import LayerPopup from "lib/popup";

interface Props {
    layerKey: string;
    type: string;
    callbackFunc?: (payload) => any;
    loading: boolean;
    nodeInfo: {
        ip: string;
        name: string;
        id: string;
    };
    allNodes: NodesData;
    nodes: NodesData;
    getAllNodes: () => any;
    createNode: (payload) => any;
    updateNode: (payload) => any;
}

const SettingNodePopup: FC<Props> = props => {
    const LanguageContext = useContext(Language());
    const { I18n } = LanguageContext;
    const [nodeInfo, setNodeInfo] = useState(
        props.type === "edit"
            ? { name: props.nodeInfo.name, ip: props.nodeInfo.ip }
            : { name: "", ip: "" }
    );

    useEffect(() => {
        props.getAllNodes();
    }, []);

    const onChangeNodeInfo = e => {
        const { name, value } = e.target;
        setNodeInfo({
            ...nodeInfo,
            [name]: value
        });
    };

    const onClickCancel = () => {
        LayerPopup.hide(props.layerKey);
    };
    const onClickAgree = () => {
        if (props.type === "create") {
            props.createNode(nodeInfo);
        } else {
            const payload = {
                ...nodeInfo,
                id: props.nodeInfo.id
            };
            props.updateNode(payload);
        }
        LayerPopup.hide(props.layerKey);
    };

    return (
        <div className="popup-content">
            <p className="title">
                {props.type === "create" ? I18n.addNode : I18n.editNode}
            </p>
            <div className="setting-form">
                <div className="setting-input">
                    <label htmlFor="nodeName">{I18n.nodeName}</label>
                    <input
                        type="text"
                        name="name"
                        onChange={onChangeNodeInfo}
                        value={nodeInfo.name}
                        autoFocus={props.type === "create" && true}
                    />
                </div>
                <div className="setting-input">
                    <label htmlFor="setting">{I18n.nodeIP}</label>
                    <input
                        type="text"
                        name="ip"
                        onChange={onChangeNodeInfo}
                        value={nodeInfo.ip}
                    />
                </div>
            </div>
            <div className="btn-box">
                <Button className="cancel-btn" onClick={onClickCancel}>
                    {I18n.cancel}
                </Button>
                <Button
                    className="agree-btn"
                    onClick={onClickAgree}
                    disabled={!nodeInfo.name || !nodeInfo.ip}
                >
                    {I18n.agree}
                </Button>
            </div>
        </div>
    );
};

const mapStateToProps = state => ({
    allNodes: state.nodes.allNodes,
    nodes: state.nodes.nodes,
    loading: state.nodes.loading
});
const mapDispatchToProps = (dispatch: Dispatch<Action>) => ({
    getAllNodes: () => dispatch(nodesActions.getAllNodes()),
    createNode: payload => dispatch(nodesActions.createNode(payload)),
    updateNode: payload => dispatch(nodesActions.updateNode(payload))
});

export default connect(
    mapStateToProps,
    mapDispatchToProps
)(SettingNodePopup);
