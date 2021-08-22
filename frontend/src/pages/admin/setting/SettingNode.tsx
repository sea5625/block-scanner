import { connect } from "react-redux";
import { Dispatch, Action } from "redux";
import { SettingNode } from "components";
import { actions as nodesActions } from "modules/nodes";

const mapStateToProps = state => ({
    allNodes: state.nodes.allNodes,
    loading: state.nodes.loading
});
const mapDispatchToProps = (dispatch: Dispatch<Action>) => ({
    getAllNodes: () => dispatch(nodesActions.getAllNodes()),
    deleteNode: payload => dispatch(nodesActions.deleteNode(payload))
});

export default connect(
    mapStateToProps,
    mapDispatchToProps
)(SettingNode);
