import { useEffect } from "react";
import PropTypes from "prop-types";

function BlackList(){

    return <>
        black list
        </>
}

BlackList.propTypes = {
    setIsLoggedIn: PropTypes.func.isRequired
}
export default BlackList;