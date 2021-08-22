import React, { useContext } from "react";
import { MdList, MdLaptop } from "react-icons/md";
import { DiNetbeans } from "react-icons/di";
import { Language } from "components";

const Total = props => {
    const LanguageContext = useContext(Language());
    const { I18n } = LanguageContext;

    return (
        <div className="total-box">
            <div className="total-content node">
                <MdLaptop />
                <p className="total-count">{props.channel.data.nodes.length}</p>
                <p className="total-title">{I18n.totalNodes}</p>
            </div>
            <div className="total-content block">
                <DiNetbeans />
                <p className="total-count">{props.channel.data.blockHeight}</p>
                <p className="total-title">{I18n.totalBlockHeight}</p>
            </div>
            <div className="total-content tx">
                <MdList />
                <p className="total-count">{props.channel.data.countOfTX}</p>
                <p className="total-title">{I18n.totalTx}</p>
            </div>
        </div>
    );
};

export default Total;
