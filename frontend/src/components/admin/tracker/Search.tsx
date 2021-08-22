import React, { FC, useState, useContext } from "react";
import { connect } from "react-redux";
import { Dispatch, Action } from "redux";
import { history } from "store";
import { actions as trackerActions } from "modules/tracker";
import { Language } from "components";

interface Props {
    getBlockOrTxBySearch: (payload) => any;
}
//[TO DO]
const Search = props => {
    const LanguageContext = useContext(Language());
    const { I18n } = LanguageContext;
    const [searchValue, setSearchValue] = useState("");

    const onKeyPress = e => {
        if (e.charCode === 13) {
            onSubmit();
        }
    };
    const onSubmit = () => {
        const { channel, name } = props;
        props.getBlockOrTxBySearch({ channel, name, searchValue });
        setSearchValue("");
    };
    return (
        <div className="search-box">
            <label>{I18n.search} :</label>
            <input
                className="search-input"
                type="text"
                onChange={e => setSearchValue(e.target.value)}
                value={searchValue}
                onKeyPress={onKeyPress}
                placeholder={I18n.searchTxOrBlockByHash}
            />
        </div>
    );
};

const mapStateToProps = state => ({
    language: state.storage.language
});

const mapDispatchToProps = (dispatch: Dispatch<Action>) => ({
    getBlockOrTxBySearch: payload =>
        dispatch(trackerActions.getBlockOrTxBySearch(payload))
});

export default connect(
    mapStateToProps,
    mapDispatchToProps
)(Search);
