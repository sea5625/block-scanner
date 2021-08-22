import { connect } from "react-redux";
import { Dispatch, Action } from "redux";
import { BlockList } from "components";
import { actions as trackerActions } from "modules/tracker";

const mapStateToProps = state => ({
    blocksLoading: state.tracker.blocksLoading,
    blocks: state.tracker.blocks
});
const mapDispatchToProps = (dispatch: Dispatch<Action>) => ({
    getBlocks: payload => dispatch(trackerActions.getBlocks(payload))
});

export default connect(
    mapStateToProps,
    mapDispatchToProps
)(BlockList);
