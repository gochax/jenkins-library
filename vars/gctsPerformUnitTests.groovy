import com.sap.piper.PiperGoUtils
import com.sap.piper.Utils
import groovy.transform.Field

import static com.sap.piper.Prerequisites.checkScript

@Field String STEP_NAME = getClass().getName()
@Field String METADATA_FILE = 'metadata/gctsRunUnitTestsForAllRepoPackages.yaml'

void call(Map parameters = [:]) {
        List credentials = [
        [type: 'usernamePassword', id: 'credentialsId', env: ['PIPER_username', 'PIPER_password']]
        ]
        piperExecuteBin(parameters, STEP_NAME, METADATA_FILE, credentials)
}