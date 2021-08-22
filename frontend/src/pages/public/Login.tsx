import { connect } from "react-redux";
import { Dispatch, Action } from "redux";
import { Login } from "components";
import { actions as authActions } from "modules/auth";
import { actions as storageActions } from "modules/storage";
import { actions as resourceActions } from "modules/resource";

const mapStateToProps = state => ({
    authStatus: state.auth.payload,
    errorMessage: state.auth.errorMessage,
    resource: state.resource.resource,
    loading: state.resource.loading
});
const mapDispatchToProps = (dispatch: Dispatch<Action>) => ({
    login: payload => dispatch(authActions.login(payload)),
    setToken: payload => dispatch(storageActions.setToken(payload)),
    getResource: payload => dispatch(resourceActions.getResource(payload))
});

export default connect(
    mapStateToProps,
    mapDispatchToProps
)(Login);
