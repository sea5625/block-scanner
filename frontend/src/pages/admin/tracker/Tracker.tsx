import { connect } from "react-redux";
import { Dispatch, Action } from "redux";
import { Tracker } from "components";
import { actions as trackerActions } from "modules/tracker";
import { actions as channelsActions } from "modules/channels";

const mapStateToProps = state => ({
    blocksLoading: state.tracker.blocksLoading,
    txsLoading: state.tracker.txsLoading,
    blocks: state.tracker.blocks,
    txs: state.tracker.txs,
    channel: state.channels.channel
});
const mapDispatchToProps = (dispatch: Dispatch<Action>) => ({
    getBlocks: payload => dispatch(trackerActions.getBlocks(payload)),
    getTxs: payload => dispatch(trackerActions.getTxs(payload)),
    getChannel: payload => dispatch(channelsActions.getChannel(payload))
});

export default connect(
    mapStateToProps,
    mapDispatchToProps
)(Tracker);
