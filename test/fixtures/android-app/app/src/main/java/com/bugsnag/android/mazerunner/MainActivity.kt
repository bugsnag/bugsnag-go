package com.bugsnag.android.mazerunner

import android.os.Build
import android.os.Bundle
import android.support.v7.app.AppCompatActivity
import android.util.Log
import com.bugsnag.android.*

class MainActivity : AppCompatActivity() {

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)
    }

    override fun onResume() {
        super.onResume()
        initialiseBugsnag()
        enqueueTestCase()
    }

    /**
     * Enqueues the test case with a delay on the main thread. This avoids the Activity wrapping
     * unhandled Exceptions
     */
    private fun enqueueTestCase() {
        window.decorView.postDelayed({
            throw RuntimeException("UnhandledExceptionScenario")
        }, 100)
    }

    private fun initialiseBugsnag() {
        val config = Configuration(intent.getStringExtra("BUGSNAG_API_KEY"))
        val port = intent.getStringExtra("BUGSNAG_PORT")
        config.endpoint = "${findHostname()}:$port"
        config.sessionEndpoint = "${findHostname()}:$port"

        Bugsnag.init(this, config)
        Bugsnag.setLoggingEnabled(true)
    }

    private fun findHostname(): String {
        val isEmulator = Build.FINGERPRINT.startsWith("unknown")
                || Build.FINGERPRINT.contains("generic")
        return when {
            isEmulator -> "http://10.0.2.2"
            else -> "http://localhost"
        }
    }

}
