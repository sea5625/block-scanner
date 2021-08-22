import React from "react";
import { connect } from "react-redux";
import AdminRoutes from "pages/admin/AdminRoutes";
import Header from "components/admin/Header";
import Sidebar from "components/admin/Sidebar";

const Admin = () => {
    return (
        <div className="wrap">
            <Header />
            <Sidebar />
            <div className="container">
                <AdminRoutes />
            </div>
        </div>
    );
};

const mapStateToProps = state => ({
    language: state.storage.language
});

export default connect(
    mapStateToProps,
    null
)(Admin);
