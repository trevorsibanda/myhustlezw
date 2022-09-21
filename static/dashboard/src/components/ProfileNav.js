import React, {Component} from "react";
import {Link} from "react-router-dom"



class ProfileNav extends Component {
    
    render() {
        return (
            <div class="row">
                <div class="col-md-4 col-xs-12" style={{ 'font-size': '18px' }}>
                    <a href={"/@" + this.props.user.username}>@<span>{this.props.user.username}</span></a>&nbsp;
                </div>
                <div class="col-md-8 col-xs-12 d-none d-md-block">
                    <div class="btn-wrapper pull-right">
                        <Link to="/content/new" class="boxed-btn btn-business "><span class="fa fa-plus"></span>
                      &nbsp;Create new campaign</Link>
                    </div>
                </div>
                <div class="col-md-8 col-xs-12 d-xs-block d-sm-block d-md-none">
                    <div class="btn-wrapper ">
                        <Link to="/creator/content/new" class="boxed-btn btn-business btn-block"><span class="fa fa-plus"></span>
                      &nbsp;Create new campaign</Link>
                    </div>
                </div>
            </div>
        );
    }
}

export default ProfileNav;