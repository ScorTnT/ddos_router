import { useEffect } from "react";
import PropTypes from "prop-types";

function BlackList(){
    
    useEffect(() => {
    }, []);
    return <>
        black list
        </>
}

BlackList.propTypes = {
    setIsLoggedIn: PropTypes.func.isRequired
}
export default BlackList;