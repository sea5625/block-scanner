import { connect } from "react-redux";
import { Dispatch, Action } from "redux";
import { NodeList } from "components";
import { actions as nodesActions } from "modules/nodes";

const mapStateToProps = state => ({
    nodes: state.nodes.nodes,
    loading: state.nodes.loading
});
const mapDispatchToProps = (dispatch: Dispatch<Action>) => ({
    getNodes: payload => dispatch(nodesActions.getNodes(payload))
});

export default connect(
    mapStateToProps,
    mapDispatchToProps
)(NodeList);
