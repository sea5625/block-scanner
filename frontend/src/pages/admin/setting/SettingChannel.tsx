import { connect } from "react-redux";
import { Dispatch, Action } from "redux";
import { SettingChannel } from "components";
import { actions as channelsActions } from "modules/channels";

const mapStateToProps = state => ({
    channels: state.channels.channels,
    loading: state.channels.loading
});
const mapDispatchToProps = (dispatch: Dispatch<Action>) => ({
    getChannels: () => dispatch(channelsActions.getChannels()),
    deleteChannel: payload => dispatch(channelsActions.deleteChannel(payload))
});

export default connect(
    mapStateToProps,
    mapDispatchToProps
)(SettingChannel);
