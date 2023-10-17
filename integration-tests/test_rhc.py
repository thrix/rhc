import logging
from packaging import version

def test_version(rhc, subtests):
    with subtests.test("RHC version"):
        proc = rhc.run("--version")
        assert b"rhc version" in proc.stdout
    
    with subtests.test("RHC version in the proper format"):
        version_string = proc.stdout.decode().strip()
        rhc_version = version.parse(version_string.split()[2])
        assert isinstance(rhc_version, version.Version) is True


def test_rhc_register(external_candlepin, test_config, rhc, subtests):
    assert not rhc.is_registered
    with subtests.test(msg="RHC connect"):
        proc= rhc.connect(
            username=test_config.get("candlepin", "username"),
            password=test_config.get("candlepin", "password"),
        )
        logging.info(f'result is {proc.stdout}')
        assert "Successfully connected to Red Hat!" in proc.stdout.decode()
        assert rhc.is_registered

    with subtests.test(msg="RHC disconnect"):
        rhc.disconnect()
        logging.info(f'result of disconnect task: {proc.stdout}')
        assert not rhc.is_registered


def test_rhc_status(rhc):
    proc = rhc.run("status")
    logging.info(f'rhc status is: {proc.stdout}')
