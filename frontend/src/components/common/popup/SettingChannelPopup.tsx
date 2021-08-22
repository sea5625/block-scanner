import React, { FC, useState, useEffect, useContext } from "react";
import { connect } from "react-redux";
import { Dispatch, Action } from "redux";
import { withStyles } from "@material-ui/core/styles";
import { Checkbox, FormControlLabel, Button } from "@material-ui/core";
import { actions as nodesActions } from "modules/nodes";
import { actions as channelsActions } from "modules/channels";
import { Language } from "components";
import { NodesData } from "modules/nodes";
import LayerPopup from "lib/popup";

interface Props {
    layerKey: string;
    type: string;
    callbackFunc?: (payload) => any;
    loading: boolean;
    channelInfo: {
        channel: string;
        name: string;
    };
    allNodes: NodesData;
    nodes: NodesData;
    getAllNodes: () => any;
    getNodes: (payload) => any;
    createChannel: (payload) => any;
    updateChannel: (payload) => any;
}
const CustomCheckbox = withStyles({
    root: {
        color: "rgb(94, 114, 228)",
        "&$checked": {
            color: "rgb(94, 114, 228)"
        },
        width: 10,
        height: 10
    },
    checked: {}
})(props => <Checkbox color="default" {...props} />);

const SettingChannelPopup: FC<Props> = props => {
    const LanguageContext = useContext(Language());
    const { I18n } = LanguageContext;
    const [channelName, setChannelName] = useState(
        props.channelInfo ? props.channelInfo.name : "default"
    );
    const [checkedArr, setCheckedArr] = useState([]);

    useEffect(() => {
        props.getAllNodes();
        if (props.type === "edit") {
            props.getNodes(props.channelInfo.channel);
        }
    }, []);

    useEffect(() => {
        if (props.type === "edit") {
            makeEditCheckedArr();
        }
    }, [props.allNodes, props.nodes]);

    const makeEditCheckedArr = () => {
        let allArr = [];
        props.allNodes.data.forEach(el => {
            allArr.push(el.name);
        });
        let defaultArr = [];
        props.nodes.data.forEach(el => {
            defaultArr.push(el.name);
        });

        const newArr = allArr.filter(el => {
            return defaultArr.includes(el);
        });
        setCheckedArr(newArr);
    };

    const onChangeCheckbox = e => {
        const { name } = e.target;
        if (checkedArr.includes(name)) {
            const newArr = checkedArr.filter(el => {
                return el !== name;
            });
            setCheckedArr(newArr);
        } else {
            setCheckedArr(checkedArr.concat(name));
        }
    };
    const onClickCancel = () => {
        LayerPopup.hide(props.layerKey);
    };
    const onClickAgree = () => {
        let nodes = [];
        props.allNodes.data.forEach(el => {
            if (checkedArr.includes(el.name)) {
                nodes.push(el.id);
            }
        });
        if (props.type === "create") {
            const payload = {
                data: {
                    name: channelName,
                    nodes
                }
            };
            props.createChannel(payload);
        } else {
            const payload = {
                data: {
                    name: channelName,
                    nodes
                },
                id: props.channelInfo.channel
            };
            props.updateChannel(payload);
        }
        LayerPopup.hide(props.layerKey);
    };
    return (
        <div className="popup-content">
            <p className="title">
                {props.type === "create" ? I18n.addChannel : I18n.editChannel}
            </p>
            <div className="setting-form">
                <div className="channel-name-input">
                    <label htmlFor="channel-name">{I18n.channelName}</label>
                    <input
                        type="text"
                        onChange={e => setChannelName(e.target.value)}
                        value={channelName}
                        autoFocus={true}
                    />
                </div>
                <div className="select-node-list">
                    <p className="list-title">{I18n.selectedNode}</p>
                    <ul>
                        {props.allNodes.data.map((el, key) => {
                            return (
                                <li key={key}>
                                    <FormControlLabel
                                        label={el.name}
                                        name={el.name}
                                        onChange={onChangeCheckbox}
                                        checked={checkedArr.includes(el.name)}
                                        control={<CustomCheckbox />}
                                    />
                                </li>
                            );
                        })}
                    </ul>
                </div>
                <div className="selected-node-list">{checkedArr.join()}</div>
            </div>
            <div className="btn-box">
                <Button className="cancel-btn" onClick={onClickCancel}>
                    {I18n.cancel}
                </Button>
                <Button
                    className="agree-btn"
                    onClick={onClickAgree}
                    disabled={checkedArr.length < 1 || !channelName}
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
    getNodes: payload => dispatch(nodesActions.getNodes(payload)),
    createChannel: payload => dispatch(channelsActions.createChannel(payload)),
    updateChannel: payload => dispatch(channelsActions.updateChannel(payload))
});

export default connect(
    mapStateToProps,
    mapDispatchToProps
)(SettingChannelPopup);
